---
title: Cursor policy gateway for MCP tools
description: Put allow, review, and block policy in front of MCP tools used from Cursor.
---

# Cursor policy gateway for MCP tools

Cursor can call MCP servers with credentials and access to external systems. AegisFlow sits between Cursor MCP client and upstream MCP server. Every routed tool call becomes an `ActionEnvelope` before execution.

## Local setup

Run PR-writer starter kit:

```bash
git clone https://github.com/saivedant169/AegisFlow.git
cd AegisFlow/starter-kit
./install-pr-writer.sh
```

Add bridge from [Cursor setup file](https://github.com/saivedant169/AegisFlow/blob/main/starter-kit/editors/cursor.md) to project MCP configuration. Bridge sends JSON-RPC requests to local AegisFlow MCP gateway.

## Policy decisions

Typical PR-writer flow:

| Cursor MCP action | Decision |
|---|---|
| Read repository metadata | allow |
| Run approved test command | allow |
| Open pull request | review |
| Delete repository | block |
| Merge pull request | block |

Review action stays pending until reviewer approves or denies it. Approval applies to one matching action and is consumed once.

## Boundary limitation

Cursor built-in file edits, terminal commands, and network calls are outside AegisFlow unless routed through configured gateway. AegisFlow is policy boundary, not desktop sandbox. Keep operating-system permissions and workspace isolation in place.
