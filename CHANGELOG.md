# Changelog

## v0.7.0 — 2026-10-30

- Refresh the official OpenAPI snapshot and complete generated client: 297 paths,
  337 operations, and 338 schemas, up from 288, 326, and 313 in v0.6.0.
  The upstream API still identifies itself as `0.9.15`.
- Add generated APIs for Identity profiles, Memory identity evidence, extraction
  and correction, entity aliases, guarded Timeline review, Thread projects, and
  JSON/multipart import previews. No upstream paths or operations were removed.
- Add concise service helpers for project listing/filtering, import previews,
  scoped FS deletion, and detailed storage optimization results.
- Preserve Memory write intent and source assessment, native source-derived
  distillation, extraction policy overrides, and distill-again task status.
- Expose Timeline inclusion options for all data transfer requests and multipart
  attachment references for file ingestion.
- Preserve Thread deletion partial/uncertain outcomes and repair flags, Agent
  identity references, and exact KG resolution identities and warnings.
- Keep existing service method signatures and legacy response fields. New complex
  types reuse the generated `openapi` package; generated types track upstream
  schema changes, including structured bulk-deletion result entries.
- Add wire-contract and JSON round-trip regression coverage. The upstream fixed
  both previously malformed schedule query parameters; older snapshots remain
  supported by the existing reviewed corrections.

The minimum Go version remains 1.27.
