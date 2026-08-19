# Coding Harness Expansion Plan

Status: proposed

Snapshot date: 2026-08-20

## Decision

AegisFlow should become a protocol-first policy boundary, not a collection of product-specific proxies.

Six integration paths cover most current coding harnesses:

1. Native pre-tool hooks for built-in file, shell, network, and MCP actions.
2. Model Context Protocol gateway for external tools.
3. Agent Client Protocol proxy for editor-mediated file, terminal, permission, and MCP calls.
4. Isolated runner for open-source command-line harnesses without a blocking hook.
5. Provider API proxy for model traffic and products with configurable provider endpoints.
6. Remote relay for hosted harnesses that accept HTTPS MCP, API, or webhook integration.

Each path has a different security boundary. Product pages must state that boundary. MCP support alone must never be described as control over native shell or file tools.

## Census Method

No authoritative registry covers every coding harness. GitHub topic search also cannot provide a valid product count. On snapshot date, `topic:coding-agent archived:false` returned 4,205 repositories. Only 244 were unarchived, had at least 100 stars, and had activity within 180 days. Most were libraries, prompts, wrappers, examples, or orchestration tools rather than runnable coding harnesses.

Primary census uses product families. A family qualifies through at least one gate:

1. Listed in official Agent Client Protocol registry.
2. Public, installable, vendor-backed coding product with official documentation.
3. Runnable open-source coding harness, unarchived, updated within 180 days, and holding at least 1,000 GitHub stars.

Duplicates across CLI, IDE, desktop, and cloud surfaces count once when they share one product family. Review-only bots, autocomplete-only products, app builders, SDKs, model providers, prompt packs, and agent orchestration frontends stay outside primary count.

Result: **77 verified product families**.

This number is reproducible for snapshot date, not permanent. Monthly refresh must update source commit, repository metadata, additions, removals, and status changes.

## Census Result

### ACP Registry: 38 Families

Official registry contained 39 manifests at commit [`8f0edb7`](https://github.com/agentclientprotocol/registry/commit/8f0edb76dff2f24f7b1912e397c1025ae9364f7b). Two manifests represent GitHub Copilot family, leaving 38 distinct families.

| # | Product family | Primary integration |
|---:|---|---|
| 1 | Agoragentic | ACP |
| 2 | Amp | ACP |
| 3 | Auggie | ACP |
| 4 | Autohand Code | ACP |
| 5 | Claude Code | Native hook, ACP, MCP |
| 6 | Cline | ACP, MCP |
| 7 | CodeBuddy Code | ACP |
| 8 | Codex | ACP, provider proxy |
| 9 | Cortex Code | ACP |
| 10 | Corust Agent | ACP |
| 11 | crow-cli | ACP |
| 12 | Cursor | Native hook, ACP, MCP |
| 13 | DeepAgents | ACP |
| 14 | Devin | Native hook, ACP, MCP |
| 15 | DimCode | ACP |
| 16 | Dirac | ACP |
| 17 | Factory Droid | ACP |
| 18 | fast-agent | ACP |
| 19 | Gemini CLI | Native hook, ACP, MCP |
| 20 | GitHub Copilot | Native hook, ACP, MCP |
| 21 | GLM Agent | ACP |
| 22 | Goose | ACP, MCP |
| 23 | Grok Build | ACP |
| 24 | Harn | ACP |
| 25 | Junie | ACP |
| 26 | Kilo | ACP, MCP |
| 27 | Kimi CLI | ACP |
| 28 | Minion Code | ACP |
| 29 | Mistral Vibe | ACP |
| 30 | Nova | ACP |
| 31 | OpenCode | ACP, MCP |
| 32 | pi | ACP |
| 33 | Poolside | ACP |
| 34 | Qoder | ACP |
| 35 | Qwen Code | ACP, MCP |
| 36 | siGit | ACP |
| 37 | Stakpak | ACP |
| 38 | VT Code | ACP |

ACP makes these products addressable only where operations cross ACP. Agent-private shell, file, browser, or network actions remain outside that boundary unless another adapter or sandbox captures them.

### Vendor Products Outside Registry: 19 Families

| # | Product family | Planned integration | Current certainty |
|---:|---|---|---|
| 1 | [Amazon Q Developer](https://docs.aws.amazon.com/amazonq/latest/qdeveloper-ug/command-line-mcp-config-CLI.html) | MCP, provider proxy | Tool-only until tested |
| 2 | [Google Antigravity](https://antigravity.google/docs/hooks) | Native hook, MCP | Pre-tool blocking documented |
| 3 | [BLACKBOX CLI](https://docs.blackbox.ai/features/blackbox-cli/getting-started) | Isolated runner, provider proxy | Runner path required |
| 4 | [GitLab Duo Agent Platform](https://docs.gitlab.com/user/duo_agent_platform/) | MCP, remote relay | MCP remains beta |
| 5 | [IBM Bob](https://bob.ibm.com/docs/ide/configuration/mcp/mcp-in-bob) | MCP | Tool-only until tested |
| 6 | [Jules](https://jules.google/docs/changelog/2026-02-02) | Partner integration | Curated MCP servers only |
| 7 | [Kiro](https://kiro.dev/docs/cli/hooks/) | Native hook, MCP | Pre-tool blocking documented |
| 8 | [Qodo](https://docs.qodo.ai/qodo-ide/tools-mcps/agentic-tools-mcps) | MCP | Tool-only until tested |
| 9 | [Replit Agent](https://docs.replit.com/build/connect-via-mcp) | Remote MCP | Custom HTTPS MCP documented |
| 10 | [Roomote](https://github.com/RooCodeInc/Roomote) | Isolated runner, remote relay | Self-hosted path available |
| 11 | [Rovo Dev CLI](https://support.atlassian.com/rovo/docs/use-rovo-dev-cli/) | Hook spike, MCP | Exact blocking contract needs test |
| 12 | [Salesforce Agentforce Vibes IDE](https://developer.salesforce.com/blogs/2026/04/new-developer-edition-agentforce-vibes-claude-mcp) | Partner integration | Product exposes MCP server, not verified custom MCP client |
| 13 | [Sourcegraph Cody Enterprise](https://sourcegraph.com/docs/cody/capabilities/agentic-context-fetching) | MCP | Local MCP tools only |
| 14 | [Tabnine Agent](https://docs.tabnine.com/main/getting-started/tabnine-agent) | Native hook, MCP | Pre-tool blocking documented |
| 15 | [Trae Agent](https://github.com/bytedance/trae-agent) | Isolated runner, provider proxy | Open-source CLI path |
| 16 | [Warp Oz](https://docs.warp.dev/reference/cli) | MCP, remote relay | Tool-only until tested |
| 17 | [Zed Agent](https://zed.dev/docs/ai/agents) | MCP | Native tools outside MCP boundary |
| 18 | [Zencoder](https://docs.zencoder.ai/features/coding-agent) | MCP | Tool-only until tested |
| 19 | [Zoo Code](https://github.com/Zoo-Code-Org/Zoo-Code) | MCP, isolated runner | Successor to Roo Code |

### Active Open-Source Harnesses: 20 Families

| # | Product family | Repository | Planned integration |
|---:|---|---|---|
| 1 | Aider | [Aider-AI/aider](https://github.com/Aider-AI/aider) | Isolated runner, provider proxy |
| 2 | OpenHands | [OpenHands/OpenHands](https://github.com/OpenHands/OpenHands) | Isolated runner, security plugin |
| 3 | Crush | [charmbracelet/crush](https://github.com/charmbracelet/crush) | MCP, isolated runner |
| 4 | Open Interpreter | [openinterpreter/openinterpreter](https://github.com/openinterpreter/openinterpreter) | MCP, isolated runner |
| 5 | Open SWE | [langchain-ai/open-swe](https://github.com/langchain-ai/open-swe) | MCP, isolated runner |
| 6 | mini-SWE-agent | [SWE-agent/mini-swe-agent](https://github.com/SWE-agent/mini-swe-agent) | Environment adapter, isolated runner |
| 7 | CodeWhale | [Hmbown/CodeWhale](https://github.com/Hmbown/CodeWhale) | Isolated runner |
| 8 | DeepSeek Reasonix | [esengine/DeepSeek-Reasonix](https://github.com/esengine/DeepSeek-Reasonix) | Isolated runner, provider proxy |
| 9 | jcode | [1jehuang/jcode](https://github.com/1jehuang/jcode) | Isolated runner |
| 10 | Kun | [KunAgent/Kun](https://github.com/KunAgent/Kun) | Isolated runner |
| 11 | Command Code | [CommandCodeAI/command-code](https://github.com/CommandCodeAI/command-code) | Isolated runner |
| 12 | grok-cli | [superagent-ai/grok-cli](https://github.com/superagent-ai/grok-cli) | Isolated runner, provider proxy |
| 13 | dao-code | [tigicion/dao-code](https://github.com/tigicion/dao-code) | MCP, native hook, isolated runner |
| 14 | little-coder | [itayinbarr/little-coder](https://github.com/itayinbarr/little-coder) | Isolated runner |
| 15 | CoreCoder | [he-yufeng/CoreCoder](https://github.com/he-yufeng/CoreCoder) | Isolated runner, provider proxy |
| 16 | zero | [Gitlawb/zero](https://github.com/Gitlawb/zero) | Isolated runner |
| 17 | zerostack | [gi-dellav/zerostack](https://github.com/gi-dellav/zerostack) | Isolated runner |
| 18 | deeptide | [paean-ai/deeptide](https://github.com/paean-ai/deeptide) | Isolated runner |
| 19 | MiniCode | [LiuMengxuan04/MiniCode](https://github.com/LiuMengxuan04/MiniCode) | Isolated runner |
| 20 | pydantic-deepagents | [vstorm-co/pydantic-deepagents](https://github.com/vstorm-co/pydantic-deepagents) | Isolated runner |

### Excluded Or Deferred

| Product or class | Reason |
|---|---|
| Roo Code | Repository archived on 2026-05-15. Zoo Code is current successor. |
| Continue | Repository states active maintenance ended. |
| SWE-agent | Project directs new users to mini-SWE-agent. |
| Plandex | Activity and hosted-service status fail current gate. |
| Void | Development paused. |
| AutoCodeRover | Activity fails 180-day gate. |
| CodeRabbit and similar products | Review-only workflow, not general coding harness. |
| Lovable, Bolt, v0, Base44, Firebase Studio | Hosted app builders, not general repository harnesses. |
| LangGraph, CrewAI, AutoGen, agent SDKs | Frameworks, not installable coding products. |
| Agent frontends and multiplexers | They launch other harnesses but do not own action execution. |
| GitHub long tail below gate | Tracked for demand, not included in coverage denominator. |

## Coverage Math

Planned self-service addressability:

| Integration path | Families addressable | Boundary |
|---|---:|---|
| ACP proxy | 38 | ACP-mediated operations only |
| Vendor hooks, MCP, remote relay, or runner | 17 | Product-specific |
| Open-source isolated runner | 20 | Process, filesystem, and network boundary |
| Partner-only path | 2 | Jules and Salesforce vendor integration required |
| Total self-service | **75 of 77, 97.4%** | At least one enforceable boundary |

This is count-based addressability, not usage-weighted market share. Comparable active-user data does not exist across vendors. No release may claim 97.4 percent full action coverage.

Three support tiers prevent misleading claims:

| Tier | Claim allowed | Required proof |
|---|---|---|
| A: full action enforcement | Native writes, commands, network calls, and MCP calls checked before execution | Real binary E2E, bypass suite, fail-closed test |
| B: bounded tool enforcement | Named protocol or tool class checked before execution | Real binary E2E for that boundary |
| C: traffic visibility | Model requests or post-action events observed | Request and evidence correlation test |

Untested products remain `planned`, never `supported`.

## Target Architecture

```text
Harness
  | native hook | MCP | ACP | provider API | isolated process
  v
Adapter boundary
  v
Local policy daemon
  | normalize | evaluate | approve | record
  v
Action target or explicit block
```

### 1. Local Policy Daemon

One long-running daemon owns policy, approval state, evidence, secrets, and adapter registration.

Requirements:

* Unix socket with mode `0600` on macOS and Linux.
* Named pipe with current-user ACL on Windows.
* Optional loopback HTTP disabled by default.
* Remote mode protected by TLS, short-lived tokens, audience checks, and tenant isolation.
* Bounded request size, concurrency, queue depth, and decision timeout.
* Graceful shutdown with in-flight decision drain.
* Durable approval and evidence writes before allow response.
* Health endpoint that does not expose policy or tenant data.

### 2. Action Envelope V2

All adapters map native events into one versioned schema.

Required fields:

* Schema version.
* Harness family, surface, version, adapter ID, and adapter version.
* Session, turn, tool call, parent call, and subagent IDs.
* Actor, tenant, repository root, working directory, and environment class.
* Operation class: file, shell, SQL, Git, HTTP, MCP, browser, secret, or provider.
* Native tool name and untouched native arguments.
* Normalized targets, command, URL, method, path, repository, branch, and capability.
* Timestamp, deadline, source process, and correlation ID.
* Sensitivity labels and redaction map.

Decision values:

* `allow`
* `block`
* `review`
* `rewrite`, only when adapter declares argument mutation support

Approval token binds exact canonical action fingerprint, session, actor, adapter, expiry, and single-use nonce. Any argument, working directory, identity, or target change invalidates approval.

Canonical JSON must follow RFC 8785 before hashing or signing. Evidence records include policy version, decision reason, approval identity, adapter response, action result, and previous-record hash.

### 3. Adapter Manifest

Each adapter ships one machine-readable manifest:

```yaml
id: cursor-hooks
harness_family: cursor
surfaces: [desktop, cli]
versions:
  tested: ["x.y.z"]
capabilities:
  pre_file_read: true
  pre_file_write: true
  pre_shell: true
  pre_mcp: true
  pre_network: false
  post_result: true
  rewrite: false
failure_mode:
  write: closed
  execute: closed
  network: closed
  read: configurable
```

Manifest also records config locations, user or system scope, timeout behavior, malformed-output behavior, test artifact digest, and certification date.

### 4. Native Hook Bridge

One command handles hook JSON over stdin and emits exact product contract over stdout:

```text
aegisflow hook evaluate --adapter <id>
```

Rules:

* Logs only to stderr.
* No plain-text stdout.
* Reject unknown schema revisions in fail-closed modes.
* Preserve native arguments byte-for-byte in evidence.
* Normalize paths after resolving repository root, symlinks, and platform rules.
* Enforce deadline shorter than harness hook timeout.
* Map daemon outage to documented adapter-specific block response.
* Install system or user scoped hooks where available. Workspace-owned hook files are not a trusted security boundary because agent may edit them.
* Generate config through `aegisctl integrate <harness> --dry-run` first. Show diff, create backup, then require explicit apply.

Initial native-hook certification set:

1. Claude Code
2. Cursor
3. GitHub Copilot and VS Code Agent
4. Gemini CLI
5. Google Antigravity
6. Kiro
7. Devin Desktop and Cascade
8. Tabnine CLI

Official docs expose pre-tool blocking for these families. Product-specific failures still require real-binary tests before Tier A status.

### 5. MCP Gateway

Current AegisFlow MCP behavior targets older protocol versions. Protocol conformance comes before broader product configs.

Required protocol matrix:

* `2026-07-28`: stateless Streamable HTTP, one POST per request, request metadata, JSON or request-scoped SSE response, no protocol session, no GET stream.
* `2025-11-25`: session-based Streamable HTTP compatibility.
* `2025-06-18`: compatibility.
* `2025-03-26`: compatibility.
* `2024-11-05`: legacy HTTP plus SSE compatibility with deprecation notice.
* stdio: UTF-8 JSON-RPC framing, stderr-only logs, cancellation, process signals, and Windows support.

Security work:

* Validate `Origin` for local HTTP and bind loopback by default.
* Add DNS rebinding tests.
* Implement OAuth 2.1 resource-server behavior for remote MCP.
* Publish RFC 9728 protected-resource metadata.
* Validate token audience, issuer, scopes, PKCE flow, and step-up scope challenges.
* Support Client ID Metadata Documents, pre-registration, and legacy Dynamic Client Registration where needed.
* Cap request, SSE event, tool result, and connection sizes.
* Preserve cancellation and deadlines through upstream call.
* Keep transport parser independent from policy engine.

Use official Go SDK after current protocol support reaches stable release. Until then, pin audited commit and run spec fixtures. Do not silently track a moving pre-release dependency.

### 6. ACP Proxy

Command form:

```text
aegisflow acp wrap -- <agent command>
```

Proxy acts as ACP agent toward editor and ACP client toward wrapped agent. It mediates:

* `session/request_permission`
* `fs/read_text_file`
* `fs/write_text_file`
* `terminal/create`
* `terminal/output`
* `terminal/wait_for_exit`
* `terminal/kill`
* `terminal/release`
* MCP server configuration passed through session setup

Tool-call update notifications provide evidence but not automatic pre-execution control. Private operations performed inside wrapped process remain invisible. Tier A requires native hook or isolated runner in addition to ACP unless packet traces prove every sensitive action crosses client RPC.

Proxy requirements:

* Preserve JSON-RPC IDs and ordering.
* Handle concurrent requests and out-of-order responses.
* Propagate cancellation and process exit.
* Reject duplicate IDs and oversized messages.
* Scrub secrets from logs without changing forwarded payload.
* Pin ACP schema version and run compatibility canary against registry updates.

### 7. Isolated Runner

Harnesses without trusted pre-tool hooks run inside containment boundary.

Linux baseline:

* Rootless OCI container.
* Read-only base image.
* Workspace copy or overlay with explicit writable paths.
* No host Docker socket.
* No host SSH agent.
* No host cloud credentials.
* Dropped capabilities and `no-new-privileges`.
* CPU, memory, process, file-size, and wall-clock limits.
* Direct egress denied. HTTP and MCP traffic exits through AegisFlow.
* Short-lived secrets injected by broker and scoped to task.
* Output diff exported for review before host apply.

macOS uses a lightweight Linux VM or audited container runtime. Windows uses WSL2 or Hyper-V isolation. Native PATH shims and Git wrappers may improve visibility but do not qualify as a security boundary.

### 8. Remote Relay

Hosted harnesses require public HTTPS endpoint.

Requirements:

* Tenant-specific endpoint and OAuth client.
* mTLS option for enterprise agents.
* Region selection and retention policy.
* Durable decision queue with idempotency key.
* Web approval flow protected against replay and confused-deputy attacks.
* Webhook signature verification.
* Per-tenant rate limits and circuit breakers.
* Zero plaintext secret logging.
* Explicit data residency and deletion behavior.

Jules remains partner-only until custom MCP servers become available or vendor approves AegisFlow. Salesforce Agentforce Vibes also remains partner-only because current public material describes Salesforce as an MCP server for external clients, not Vibes as a custom MCP client. No workaround should claim native action control.

## Policy Semantics

Policy must operate on normalized capability, not harness-specific tool names.

Examples:

| Native event | Normalized capability |
|---|---|
| Shell command writing file | `filesystem.write` and `process.execute` |
| Git push | `git.remote.write` and `network.egress` |
| SQL `SELECT` | `database.read` |
| SQL `DROP TABLE` | `database.schema.delete` |
| HTTP `POST` | `network.http.write` |
| MCP GitHub create issue | `github.issue.write` |
| Secret read | `secret.read` |

Parser rules:

* Prefer structured native fields.
* Use shell parser appropriate to Bash, Zsh, PowerShell, or CMD.
* Preserve original input beside normalized form.
* Treat parse uncertainty as review or block for write and execute actions.
* Never approve by substring matching alone.

## Threat Model

Support claims must name attacker model.

| Boundary | Cooperative harness bug | Prompt injection | Compromised harness process | Malicious local user |
|---|---:|---:|---:|---:|
| MCP gateway | Yes, for MCP calls | Yes, for MCP calls | No private-action control | No |
| Native hook | Yes | Yes | Limited, process may bypass hook | No |
| ACP proxy | Yes, for ACP calls | Yes, for ACP calls | No private-action control | No |
| Provider proxy | Model traffic only | Detection only | No tool control | No |
| Isolated runner | Yes | Yes | Strong process containment goal | No host-admin protection |

Local admin remains outside threat model. Security documentation must state this.

## Test Program

### Adapter Contract Tests

Every adapter must pass:

* Allow, block, review, and supported rewrite.
* Exact action resume after approval.
* Approval replay rejection.
* Argument mutation after approval.
* Working-directory and repository mutation.
* Duplicate and parallel tool calls.
* Parent and subagent correlation.
* Daemon restart during pending review.
* Hook timeout, crash, malformed JSON, empty output, and oversized output.
* Unicode path, symlink escape, traversal, case-folding, and Windows path tests.
* Secret redaction without evidence corruption.
* Evidence tamper detection.

### MCP Conformance Tests

* One suite per supported protocol version.
* JSON and SSE response modes.
* Notification acceptance and HTTP status.
* Cancellation by stream close where required.
* Legacy session lifecycle.
* stdio framing and process termination.
* OAuth discovery, registration, PKCE, refresh, audience, issuer, and scope escalation.
* Origin and DNS rebinding attacks.
* Invalid metadata headers and mixed protocol versions.
* Slow upstream, disconnect, retry, and duplicate request.

### Real Harness E2E

Mocks cover parser branches only. Certification requires released product binary where licensing permits.

For each certified family:

1. Start clean workspace and daemon.
2. Install adapter through documented command.
3. Run allowed read.
4. Run blocked write.
5. Run blocked destructive command.
6. Run approval flow and resume exact call.
7. Attempt changed arguments with old approval.
8. Restart daemon and harness.
9. Verify result, evidence signature, and correlation IDs.
10. Remove adapter and confirm clean rollback.

Platforms:

* Linux amd64 and arm64.
* macOS amd64 and arm64.
* Windows amd64.

CI tracks one pinned stable version. Nightly canary tracks latest version. Certification expires after 30 days without successful latest-version canary.

### Release Gates

No adapter becomes `supported` until all gates pass:

* Zero bypass in declared boundary test suite.
* 100 percent decision-to-result evidence correlation.
* Approval replay and mutation tests pass.
* Fail behavior matches manifest for timeout, crash, and daemon outage.
* No secrets in stdout, logs, metrics, traces, or evidence exports.
* Local policy decision latency target: p95 at most 75 ms, p99 at most 150 ms, excluding human review.
* Race detector, fuzz tests, and parser property tests pass.
* Binary, SBOM, provenance, and checksum verification pass.
* Documentation names exact supported surface and known gaps.

## Delivery Plan

Assumption: two experienced Go engineers plus one test and release engineer. Expected elapsed time: 22 to 26 weeks. One maintainer should plan 36 to 44 weeks.

| Milestone | Duration | Deliverable | Exit condition |
|---|---:|---|---|
| H0: contract freeze | 1 week | Census, threat model, support tiers, Action Envelope V2 RFC | Maintainer approval |
| H1: adapter core | 3 weeks | Local socket API, hook bridge, manifest schema, config installer | Contract and failure tests green |
| H2: MCP current spec | 4 weeks | 2026 transport, legacy matrix, Go stdio bridge, OAuth | Protocol matrix green |
| H3: first native hooks | 4 weeks | Eight hook adapters | Real-binary E2E green on supported platforms |
| H4: ACP proxy | 4 weeks | Bidirectional proxy and registry canary | Ten registry agents tested, boundaries documented |
| H5: isolated runner | 4 weeks | Linux runner, secret broker, egress gateway | Six open-source harnesses pass bypass suite |
| H6: hosted relay | 3 weeks | Remote endpoint, OAuth, approvals, tenant isolation | Two hosted pilots pass security review |
| H7: compatibility lab | Ongoing | Nightly latest-version matrix and public status page | Monthly census and release cadence |

Parallel work allowed after H1. H3 and H4 can run together. H5 can start once Action Envelope V2 and evidence contracts freeze.

## First Certified Release

Scope should stay narrow enough to prove claims:

* Current and four legacy MCP protocol versions.
* Eight native-hook families.
* Ten ACP registry agents at Tier B.
* Six isolated open-source harnesses at Tier A containment boundary.
* Linux and macOS first. Windows certification follows after named-pipe, PowerShell, and path suites pass.

Candidate isolated set:

1. Aider
2. OpenHands
3. Crush
4. Open Interpreter
5. Open SWE
6. mini-SWE-agent

Do not add more adapters before this release passes gates. Compatibility count without proof will damage trust.

## Repository Deliverables

Planned files and ownership:

```text
schemas/action-envelope-v2.schema.json
schemas/adapter-manifest.schema.json
internal/action/
internal/adapter/
internal/hook/
internal/acp/
internal/sandbox/
internal/relay/
cmd/aegisflow-hook/
cmd/aegisflow-stdio/
cmd/aegisflow-acp/
cmd/aegisctl/
adapters/<family>/manifest.yaml
test/harnesses/<family>/
docs/compatibility/harnesses.yaml
docs/compatibility/<family>.md
```

`docs/compatibility/harnesses.yaml` becomes source for README support table and public status page. Hand-edited duplicate tables are prohibited.

## Metrics

Engineering metrics:

* Certified Tier A, B, and C family counts.
* Latest-version canary pass rate.
* Median days from harness release to adapter certification.
* Policy decision latency by adapter.
* Hook failure and fail-closed counts.
* Approval completion, expiry, and replay rejection counts.
* Evidence verification failure count.

Adoption metrics:

* Successful integrations by family.
* Weekly active protected repositories.
* Completed governed actions.
* Adapter install-to-first-action conversion.
* Adapter uninstall reasons.

Stars, page views, and raw install counts stay secondary. Successful protected workflows matter more.

## Stop Conditions

Pause expansion and fix core when any condition occurs:

* Confirmed bypass in supported boundary.
* Evidence cannot match decision to result.
* Protocol drift breaks more than 5 percent of nightly matrix.
* Local p99 decision latency exceeds 150 ms for seven days.
* More than two adapters require incompatible core semantics.
* Documentation claims exceed tested boundary.

## Source Set

Protocol and registry:

* [MCP 2026-07-28 Streamable HTTP](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/main/docs/specification/2026-07-28/basic/transports/streamable-http.mdx)
* [MCP 2026-07-28 authorization](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/main/docs/specification/2026-07-28/basic/authorization/index.mdx)
* [ACP schema](https://github.com/agentclientprotocol/agent-client-protocol/blob/main/schema/v1/schema.json)
* [ACP registry](https://github.com/agentclientprotocol/registry)
* [Zed external agents](https://zed.dev/docs/ai/external-agents)

Verified hook surfaces:

* [Claude Code hooks](https://code.claude.com/docs/en/hooks)
* [Cursor hooks](https://prod.cursor.com/docs/hooks)
* [Gemini CLI hooks](https://geminicli.com/docs/hooks/reference/)
* [VS Code agent hooks](https://code.visualstudio.com/docs/agent-customization/hooks)
* [Kiro hooks](https://kiro.dev/docs/cli/hooks/)
* [Antigravity hooks](https://antigravity.google/docs/hooks)
* [Cascade hooks](https://docs.devin.ai/desktop/cascade/hooks)
* [Tabnine hooks](https://docs.tabnine.com/main/getting-started/tabnine-cli/features/hooks)

Hosted and vendor boundaries:

* [Jules MCP policy](https://jules.google/docs/changelog/2026-02-02)
* [Replit custom MCP](https://docs.replit.com/build/connect-via-mcp)
* [GitHub cloud agent MCP](https://docs.github.com/en/copilot/concepts/agents/cloud-agent/mcp-and-cloud-agent)
* [Rovo Dev CLI](https://support.atlassian.com/rovo/docs/use-rovo-dev-cli/)
* [Sourcegraph Cody agentic context](https://sourcegraph.com/docs/cody/capabilities/agentic-context-fetching)

GitHub discovery:

* [All unarchived coding-agent topic repositories](https://github.com/search?q=topic%3Acoding-agent+archived%3Afalse&type=repositories)
* [Active repositories with at least 100 stars](https://github.com/search?q=topic%3Acoding-agent+archived%3Afalse+pushed%3A%3E%3D2026-02-20+stars%3A%3E%3D100&type=repositories)
