# QuantumNous/new-api Sync Audit

## Baselines

- Local protected baseline: `8205ec0` (2026-07-31).
- QuantumNous upstream audit head: `f1164142` (2026-08-18).
- Local worktree branch: `codex/quantum-sync`.
- The current main worktree has an unrelated, uncommitted `.gitignore` change and is not modified by this worktree.

## Protected local behavior

- Retryable 400 status codes and response retry keywords.
- Xiaomi channel integration.
- Dashboard month/yesterday/previous-day ranges and channel-affinity cache clearing.
- Usage-log channel filters, status display, usage ranges, and cache-hit statistics.
- Multi-key channels, notes, error summaries, restore/manage/test flows, and auto-disabled-key behavior.
- Codex multi-key validation/import and API-key random naming.
- Local model-mapping synchronization and pricing-precision fixes.

## Status vocabulary

- **Present**: local behavior already covers the upstream change; add only regression coverage when useful.
- **Adapt**: behavior is missing or incomplete and must be ported into local architecture.
- **Conflict**: upstream touches protected local behavior; use a three-way/manual merge and preserve local semantics.
- **Deferred**: dependency, branding, directory-layout, or large refactor changes requiring a separate compatibility decision.

## High-priority upstream queue

| Area | Upstream commits | Initial status |
| --- | --- | --- |
| Tiered billing and retry settlement | `df43f801`, `cfaba1dd` | Adapt/Conflict |
| Atomic top-up, quota concurrency, refunds | `50e5377e`, `ccd535ef`, `58d4e9bd`, `2a0ce347`, `47ba9d2c` | Adapt |
| Relay request replay and body preservation | `d6b5ce99`, `ea4f0210` | Adapt/Conflict |
| Provider conversion correctness | `2399de97`, `8ad159a3`, `253a74dd`, `7d09c695`, `4442bb30`, `3dda1d50`, `3d5dc36f`, `93d2df85` | Adapt |
| Authentication and rate-limit hardening | `d7992672`, `1da23d6b`, `9c97e78a`, `ffeb1b24`, `b6b97a66` | Adapt |
| Channel/editor capabilities | `a6cf42c0`, `57746fc9`, `e90a7c48`, `2b0efd84`, `4add708e` | Present/Adapt/Conflict |
| Frontend default adaptations | `15cfdedd`, `4eaeefbd`, `137d1171` | Adapt |

The full upstream history remains an audit source, but only entries that can be mapped without replacing protected local code are implemented on this branch.

## Verification record

The pre-change Go baseline currently has existing failures in the Claude file-content tests, stream-scanner status test, and the root package embed because `web/classic/dist` is absent. These failures are tracked separately from synchronization changes.
