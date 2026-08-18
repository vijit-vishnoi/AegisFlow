---
title: Tamper-evident agent audit logs
description: Record and verify signed, hash-linked evidence for governed agent sessions.
---

# Tamper-evident agent audit logs

AegisFlow evidence records what policy decided at configured boundary. Each session uses separate hash chain. Every record includes previous hash and HMAC signature.

## Stable signing key

Set evidence key before starting gateway:

```bash
export AEGISFLOW_EVIDENCE_KEY=<secret-from-key-manager>
./bin/aegisflow --config configs/aegisflow.example.yaml
```

Without stable key, AegisFlow creates ephemeral key. Records from previous process cannot be verified after restart with new key.

## Verify session

```bash
./bin/aegisctl evidence export <session-id> --file evidence.json
./bin/aegisctl verify
```

Verification checks chain order, hashes, and signatures. Any edit, deletion, or reorder after signed record breaks verification from that point.

![Verified AegisFlow evidence chain](../assets/shot-evidence-verification.png)

## Security limits

Hash chain detects changes to recorded data. It does not prove calls outside AegisFlow never happened. An attacker with signing key can create valid records. Protect key outside repository, restrict admin export access, and ship evidence to append-only storage when retention matters.

See [threat model](../security/THREAT_MODEL.md) for trust assumptions.
