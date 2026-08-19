# Launch asset record

Assets captured on 2026-08-18 from local v0.9.0 build.

## Proof GIF

File: `docs/assets/hero-pr-writer.gif`

- Duration: 15.23 seconds
- Size: 1149 x 731
- Source: `scripts/record-proof.sh` recorded with asciinema, rendered with agg
- Path: real AegisFlow MCP and admin endpoints against local mock GitHub upstream
- Decisions: allow, block, review, approve, exact retry, evidence verify

## Approval replay GIF

File: `docs/assets/approval-security-e2e.gif`

- Duration: 9.08 seconds
- Size: 1084 x 706
- Source: `scripts/e2e_approval_security.sh` recorded by `scripts/record-approval-security.sh`
- Path: real binary restart, SQLite restore, exact retry, replay rejection, argument change, evidence verification, altered approval and evidence rejection

## Screenshots

| File | Content | Status |
|---|---|---|
| `shot-blocked-action.png` | `github.delete_repo` blocked with `-32001` | complete |
| `shot-approval-queue.png` | dashboard holding one `github.create_pull_request` | complete |
| `shot-evidence-verification.png` | signature and chain verification returns valid | complete |

## Benchmark card

File: `docs/assets/benchmark-card.png`

- Size: 1200 x 630
- Data: Apple M1 results from 2026-08-18
- Reproduction: `docs/performance.md`
- Editable source: `docs/assets/benchmark-card.html`

## Social preview

File: `docs/assets/social-preview.png`

- Size: 1280 x 640
- Includes product name, boundary description, decision labels, and real approval queue
- Editable source: `docs/assets/social-preview.html`

Upload social preview through repository Settings after push. GitHub repository API does not expose social preview upload.
