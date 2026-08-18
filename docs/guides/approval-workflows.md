---
title: AI agent approval workflows
description: Hold risky agent actions for human review, then authorize one exact retry.
---

# AI agent approval workflows

AegisFlow uses review decision for actions that should not run automatically but do not need permanent block. Pull request creation and database writes are common examples.

## Flow

1. Client sends routed tool call.
2. Policy returns `review`.
3. AegisFlow stores pending approval and returns structured review-required error.
4. Reviewer approves or denies pending item.
5. Client retries same action after approval.

Approval matches actor, task, protocol, tool, target, arguments, and requested capability. Request ID and timestamp may change on retry. Changed arguments require new review. Approved item is consumed after one matching call.

## CLI review

```bash
./bin/aegisctl pending
./bin/aegisctl approve <action-id> "diff and target checked"
```

Dashboard approval queue runs on admin service. Protect admin service with network controls and authentication. Reviewer identity and comment enter session evidence.

![AegisFlow approval queue](../assets/shot-approval-queue.png)

## Failure handling

Expired or denied approval does not authorize retry. If upstream call fails after approved action is consumed, operator should inspect evidence and decide whether new review is required. Do not loop retries around write actions without policy-aware error handling.
