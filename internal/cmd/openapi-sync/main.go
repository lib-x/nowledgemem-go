// Command openapi-sync validates a Nowledge Mem OpenAPI snapshot and regenerates
// the complete Go client. With -update it first refreshes the snapshot from the
// official API documentation.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

const (
	defaultOpenAPIURL = "https://mem.nowledge.co/docs/refs/openapi.json"
	generatorVersion  = "v2.8.0"
	maxSpecSize       = 16 << 20
)

type parameterCorrection struct {
	method string
	path   string
	name   string
}

var (
	pathParameterPattern = regexp.MustCompile(`\{([^}]+)\}`)
	parameterCorrections = map[parameterCorrection]string{
		{method: http.MethodGet, path: "/agent/ai-now/schedules", name: "include_deleted"}:          "query",
		{method: http.MethodGet, path: "/agent/ai-now/schedules/{schedule_id}/runs", name: "limit"}: "query",
	}
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "openapi-sync:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		sourceURL = flag.String("url", defaultOpenAPIURL, "upstream OpenAPI document URL")
		specPath  = flag.String("spec", "openapi.json", "path for the upstream OpenAPI snapshot")
		output    = flag.String("output", "client.gen.go", "path for generated Go code")
		pkg       = flag.String("package", "openapi", "generated Go package name")
		update    = flag.Bool("update", false, "download the latest upstream snapshot before generating")
	)
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var raw []byte
	var err error
	if *update {
		raw, err = download(ctx, *sourceURL)
	} else {
		raw, err = os.ReadFile(*specPath)
		if err != nil {
			err = fmt.Errorf("read OpenAPI snapshot: %w", err)
		}
	}
	if err != nil {
		return err
	}

	snapshot, err := formatJSON(raw)
	if err != nil {
		return err
	}
	prepared, version, corrections, err := prepare(snapshot)
	if err != nil {
		return err
	}

	outputDir := filepath.Dir(*output)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	generated, err := os.CreateTemp(outputDir, ".client.gen-*.go")
	if err != nil {
		return fmt.Errorf("create generated output: %w", err)
	}
	generatedPath := generated.Name()
	if err := generated.Close(); err != nil {
		return fmt.Errorf("close generated output: %w", err)
	}
	defer os.Remove(generatedPath)

	preparedFile, err := os.CreateTemp("", "nowledge-openapi-*.json")
	if err != nil {
		return fmt.Errorf("create prepared spec: %w", err)
	}
	preparedPath := preparedFile.Name()
	defer os.Remove(preparedPath)
	if _, err := preparedFile.Write(prepared); err != nil {
		preparedFile.Close()
		return fmt.Errorf("write prepared spec: %w", err)
	}
	if err := preparedFile.Close(); err != nil {
		return fmt.Errorf("close prepared spec: %w", err)
	}

	tool := "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@" + generatorVersion
	cmd := exec.CommandContext(ctx, "go", "run", tool,
		"-generate", "types,client",
		"-package", *pkg,
		"-o", generatedPath,
		preparedPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generate client: %w", err)
	}

	generatedData, err := os.ReadFile(generatedPath)
	if err != nil {
		return fmt.Errorf("read generated client: %w", err)
	}
	if err := atomicWrite(*output, generatedData, 0o644); err != nil {
		return fmt.Errorf("write generated client: %w", err)
	}
	if *update {
		if err := atomicWrite(*specPath, snapshot, 0o644); err != nil {
			return fmt.Errorf("write upstream spec: %w", err)
		}
	}

	fmt.Printf("generated %s from Nowledge Mem OpenAPI %s (%d reviewed upstream corrections)\n", *output, version, corrections)
	return nil
}

func formatJSON(raw []byte) ([]byte, error) {
	var document any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode OpenAPI snapshot: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return nil, fmt.Errorf("decode OpenAPI snapshot: %w", err)
	}

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(document); err != nil {
		return nil, fmt.Errorf("encode OpenAPI snapshot: %w", err)
	}
	return output.Bytes(), nil
}

func requireJSONEOF(decoder *json.Decoder) error {
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values")
		}
		return err
	}
	return nil
}

func download(ctx context.Context, sourceURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create OpenAPI request: %w", err)
	}
	req.Header.Set("User-Agent", "nowledgemem-go-openapi-sync")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download OpenAPI document: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download OpenAPI document: %s", resp.Status)
	}

	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxSpecSize+1))
	if err != nil {
		return nil, fmt.Errorf("read OpenAPI document: %w", err)
	}
	if len(raw) > maxSpecSize {
		return nil, fmt.Errorf("OpenAPI document exceeds %d bytes", maxSpecSize)
	}
	return raw, nil
}

func prepare(raw []byte) ([]byte, string, int, error) {
	var document map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&document); err != nil {
		return nil, "", 0, fmt.Errorf("decode OpenAPI document: %w", err)
	}
	if err := requireJSONEOF(decoder); err != nil {
		return nil, "", 0, fmt.Errorf("decode OpenAPI document: %w", err)
	}

	info, ok := document["info"].(map[string]any)
	if !ok {
		return nil, "", 0, errors.New("OpenAPI document has no info object")
	}
	if title, _ := info["title"].(string); title != "Nowledge Mem API" {
		return nil, "", 0, fmt.Errorf("unexpected OpenAPI title %q", title)
	}
	version, ok := info["version"].(string)
	if !ok || version == "" {
		return nil, "", 0, errors.New("OpenAPI document has no version")
	}
	if err := validateLocalReferences(document); err != nil {
		return nil, "", 0, err
	}

	corrections, err := correctPathParameters(document)
	if err != nil {
		return nil, "", 0, err
	}

	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(document); err != nil {
		return nil, "", 0, fmt.Errorf("encode prepared OpenAPI document: %w", err)
	}
	return output.Bytes(), version, corrections, nil
}

func validateLocalReferences(value any) error {
	switch value := value.(type) {
	case []any:
		for _, item := range value {
			if err := validateLocalReferences(item); err != nil {
				return err
			}
		}
	case map[string]any:
		for key, item := range value {
			if key == "$ref" {
				ref, ok := item.(string)
				if !ok || !strings.HasPrefix(ref, "#/") {
					return fmt.Errorf("external OpenAPI reference is not allowed: %v", item)
				}
			}
			if err := validateLocalReferences(item); err != nil {
				return err
			}
		}
	}
	return nil
}

func correctPathParameters(document map[string]any) (int, error) {
	paths, _ := document["paths"].(map[string]any)
	corrections := 0
	for path, itemValue := range paths {
		item, _ := itemValue.(map[string]any)
		expected := pathParameterNames(path)
		for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
			operation, _ := item[strings.ToLower(method)].(map[string]any)
			if operation == nil {
				continue
			}
			parameters, _ := operation["parameters"].([]any)
			declared := make([]string, 0, len(parameters))
			for _, parameterValue := range parameters {
				parameter, _ := parameterValue.(map[string]any)
				if parameter["in"] != "path" {
					continue
				}
				name, _ := parameter["name"].(string)
				if slices.Contains(expected, name) {
					declared = append(declared, name)
					continue
				}

				location, ok := parameterCorrections[parameterCorrection{method: method, path: path, name: name}]
				if !ok {
					return 0, fmt.Errorf("unreviewed path parameter mismatch: %s %s parameter %q", method, path, name)
				}
				parameter["in"] = location
				parameter["required"] = false
				corrections++
			}
			for _, name := range expected {
				if !slices.Contains(declared, name) {
					return 0, fmt.Errorf("missing path parameter: %s %s parameter %q", method, path, name)
				}
			}
		}
	}
	return corrections, nil
}

func pathParameterNames(path string) []string {
	matches := pathParameterPattern.FindAllStringSubmatch(path, -1)
	names := make([]string, 0, len(matches))
	for _, match := range matches {
		names = append(names, match[1])
	}
	return names
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".openapi-sync-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Chmod(mode); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
