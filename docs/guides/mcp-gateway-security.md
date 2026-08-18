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

## Authenticated HTTP upstreams

Keep bearer tokens out of YAML by naming their environment variable:

```yaml
mcp_gateway:
  enabled: true
  port: 8082
  upstreams:
    - name: "github"
      url: "http://127.0.0.1:8083/mcp"
      tools: ["*"]
      bearer_token_env: "GITHUB_TOKEN"
```

AegisFlow reads `GITHUB_TOKEN` at request time, sends it as a bearer token, and accepts JSON or Streamable HTTP SSE responses. Start upstream server separately and route client only through AegisFlow.

## Deployment boundary

Bind MCP gateway to loopback during local use. For remote access, place authenticated TLS proxy in front and restrict admin API separately. A client that connects directly to upstream MCP server bypasses AegisFlow.

Continue with [approval workflows](approval-workflows.md) and [production checklist](../production-checklist.md).
