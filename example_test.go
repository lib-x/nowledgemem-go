package onledgemem_test

import (
	"context"
	"fmt"
	"log"
	"time"

	mem "github.com/lib-x/nowledgemem-go"
)

func ExampleNewClient() {
	// Create a client with default settings (http://127.0.0.1:14242)
	client := mem.NewClient()
	defer client.Close()

	ctx := context.Background()

	// Health check
	health, err := client.Health.Check(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Status:", health.Status)

	// List memories
	resp, err := client.Memories.List(ctx, &mem.ListMemoriesParams{
		Limit: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range resp.Memories {
		fmt.Printf("- %s: %s\n", m.ID, m.Title)
	}
}

func ExampleNewClient_withOptions() {
	// Create a client with custom base URL and timeout
	client := mem.NewClient(
		mem.WithBaseURL("http://192.168.1.100:14242"),
		mem.WithTimeout(60*time.Second),
	)
	defer client.Close()

	ctx := context.Background()

	// Create a memory
	created, err := client.Memories.Create(ctx, &mem.CreateMemoryRequest{
		Content: "Remember to review the API design",
		Title:   strPtr("API Review"),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created memory:", created.Memory.ID)

	// Search memories
	results, err := client.Memories.Search(ctx, &mem.SearchMemoriesRequest{
		Query: "API design",
		Limit: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range results.Results {
		fmt.Printf("- %s (score: %.2f)\n", r.Title, r.Score)
	}
}

func strPtr(s string) *string { return &s }
