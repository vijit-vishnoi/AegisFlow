---
title: MCP gateway security
description: Govern MCP tool discovery and calls with explicit allow, review, and block policy.
---

# MCP gateway security

Model Context Protocol connects agents to tools. Tool server may expose repository, database, shell, or SaaS operations. AegisFlow MCP gateway evaluates routed discovery and call requests before upstream execution.

## Requests under policy

Gateway handles:

- `initialize`
- `tools/list`
- `tools/call`
- supported SSE sessions

Tool call policy can match protocol, tool name, target, actor, task, requested capability, and arguments. Target glob uses doublestar semantics, so `**` crosses path separators.

## Decision behavior

`allow` forwards exact call to configured upstream. `block` returns JSON-RPC error `-32001`. `review` returns `-32002` and creates pending approval. After approval, client retries exact action. Approval matches stable action fields, ignores new request ID and timestamp, and is consumed once.

## Tool manifests

Manifest drift detection records tool inventory changes. Policy can also inspect `tools/list` responses. Review manifest changes before allowing new write-capable tools.

## Deployment boundary

Bind MCP gateway to loopback during local use. For remote access, place authenticated TLS proxy in front and restrict admin API separately. A client that connects directly to upstream MCP server bypasses AegisFlow.

Continue with [approval workflows](approval-workflows.md) and [production checklist](../production-checklist.md).
