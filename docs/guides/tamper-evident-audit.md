---
title: Tamper-evident agent audit logs
description: Record and verify signed, hash-linked evidence for governed agent sessions.
---

# Tamper-evident agent audit logs

AegisFlow evidence records what policy decided at configured boundary. Each session uses separate hash chain. Every record includes previous hash and HMAC signature.

## Restart-safe state

Memory storage remains default. Enable SQLite to restore approval and evidence records after restart:

```yaml
state:
  enabled: true
  sqlite_path: "data/aegisflow.db"
```

Set evidence key outside database before starting gateway:

```bash
export AEGISFLOW_EVIDENCE_KEY=<secret-from-key-manager>
./bin/aegisflow --config configs/aegisflow.example.yaml
```

Persistent mode refuses to start without `AEGISFLOW_EVIDENCE_KEY`. Same key authenticates approval rows and evidence records. Invalid approval signature, evidence hash, chain link, record index, or signature stops startup. `AEGISFLOW_STATE_DB` overrides configured path.

SQLite supports one running instance. Export records or ship checkpoints to separate append-only storage for longer retention and rollback detection.

## Verify session

```bash
./bin/aegisctl evidence export <session-id> --file evidence.json
./bin/aegisctl verify --session <session-id>
```

Verification checks record order, hashes, links, and signatures. Record edits, reorder, and missing records before chain tail fail verification. Tail rollback needs checkpoint stored outside same database.

![Verified AegisFlow evidence chain](../assets/shot-evidence-verification.png)

## Security limits

Hash chain detects changes to recorded data. It does not prove calls outside AegisFlow never happened. An attacker with signing key can create valid records. Protect key outside repository, restrict admin export access, and ship evidence to append-only storage when retention matters.

Run [approval replay security test](../benchmarks/approval-security.md) for restart and tamper proof.

See [threat model](../security/THREAT_MODEL.md) for trust assumptions.
