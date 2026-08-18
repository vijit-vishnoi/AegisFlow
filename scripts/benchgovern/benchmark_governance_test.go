package benchgovern

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saivedant169/AegisFlow/internal/approval"
	"github.com/saivedant169/AegisFlow/internal/credential"
	"github.com/saivedant169/AegisFlow/internal/envelope"
	"github.com/saivedant169/AegisFlow/internal/evidence"
	"github.com/saivedant169/AegisFlow/internal/toolpolicy"
)

// buildBenchEngine creates a policy engine with 20 realistic rules.
func buildBenchEngine() *toolpolicy.Engine {
	rules := make([]toolpolicy.ToolRule, 0, 20)
	protocols := []string{"mcp", "http", "shell", "sql", "git"}
	for i := 0; i < 18; i++ {
		rules = append(rules, toolpolicy.ToolRule{
			Protocol:   protocols[i%len(protocols)],
			Tool:       fmt.Sprintf("decoy.tool_%d", i),
			Decision:   "block",
			Capability: "delete",
		})
	}
	rules = append(rules, toolpolicy.ToolRule{
		Protocol: "shell",
		Tool:     "shell.exec",
		Target:   "terraform*",
		Decision: "review",
	})
	rules = append(rules, toolpolicy.ToolRule{
		Protocol: "mcp",
		Tool:     "github.list_*",
		Decision: "allow",
	})
	return toolpolicy.NewEngine(rules, "block")
}

func benchEnvelope() *envelope.ActionEnvelope {
	return envelope.NewEnvelope(
		envelope.ActorInfo{
			Type:      "agent",
			ID:        "bench-agent",
			SessionID: "bench-session",
			TenantID:  "bench-tenant",
		},
		"benchmark-task",
		envelope.ProtocolMCP,
		"github.list_repos",
		"github.com/org",
		envelope.CapRead,
	)
}

func benchBlockEnvelope() *envelope.ActionEnvelope {
	return envelope.NewEnvelope(
		envelope.ActorInfo{
			Type:      "agent",
			ID:        "bench-agent",
			SessionID: "bench-session",
			TenantID:  "bench-tenant",
		},
		"benchmark-task",
		envelope.ProtocolHTTP,
		"unknown.tool",
		"some-target",
		envelope.CapWrite,
	)
}

func benchReviewEnvelope() *envelope.ActionEnvelope {
	return envelope.NewEnvelope(
		envelope.ActorInfo{
			Type:      "agent",
			ID:        "bench-agent",
			SessionID: "bench-session",
			TenantID:  "bench-tenant",
		},
		"benchmark-task",
		envelope.ProtocolShell,
		"shell.exec",
		"terraform apply",
		envelope.CapDeploy,
	)
}

type benchCredentialBroker struct{}

func (benchCredentialBroker) Issue(_ context.Context, req credential.CredentialRequest) (*credential.Credential, error) {
	issuedAt := time.Now()
	return &credential.Credential{
		ID:        "bench-credential",
		Type:      "benchmark",
		Token:     "bench-token",
		ExpiresAt: issuedAt.Add(req.TTL),
		Scope:     req.Capability + ":" + req.Target,
		TaskID:    req.TaskID,
		IssuedAt:  issuedAt,
	}, nil
}

func (benchCredentialBroker) Revoke(context.Context, string) error { return nil }
func (benchCredentialBroker) Name() string                         { return "benchmark" }

func BenchmarkPolicyEvaluateAllow(b *testing.B) {
	engine := buildBenchEngine()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := benchEnvelope()
		engine.Evaluate(env)
	}
}

func BenchmarkPolicyEvaluateBlock(b *testing.B) {
	engine := buildBenchEngine()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := benchBlockEnvelope()
		engine.Evaluate(env)
	}
}

func BenchmarkPolicyEvaluateWithEvidence(b *testing.B) {
	engine := buildBenchEngine()
	chain := evidence.NewSessionChain("bench-session")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := benchEnvelope()
		d := engine.Evaluate(env)
		env.PolicyDecision = d
		chain.Record(env)
	}
}

func BenchmarkFullAllowPipeline(b *testing.B) {
	engine := buildBenchEngine()
	chain := evidence.NewSessionChain("bench-session")
	reg := credential.NewRegistry()
	reg.Register("benchmark", benchCredentialBroker{})
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := benchEnvelope()
		d := engine.Evaluate(env)
		env.PolicyDecision = d
		chain.Record(env)
		reg.Issue(ctx, "benchmark", credential.CredentialRequest{
			TaskID:     "bench-task",
			SessionID:  "bench-session",
			TenantID:   "bench-tenant",
			Tool:       env.Tool,
			Target:     env.Target,
			Capability: string(env.RequestedCapability),
			TTL:        5 * time.Minute,
		}, env.ID)
	}
}

func BenchmarkReviewPipeline(b *testing.B) {
	engine := buildBenchEngine()
	queue := approval.NewQueue(b.N + 1000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := benchReviewEnvelope()
		d := engine.Evaluate(env)
		env.PolicyDecision = d
		queue.Submit(env)
	}
}

func BenchmarkEvidenceRecordOnly(b *testing.B) {
	chain := evidence.NewSessionChain("bench-session")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env := benchEnvelope()
		env.PolicyDecision = envelope.DecisionAllow
		chain.Record(env)
	}
}

func BenchmarkEnvelopeCreation(b *testing.B) {
	actor := envelope.ActorInfo{
		Type:      "agent",
		ID:        "bench-agent",
		SessionID: "bench-session",
		TenantID:  "bench-tenant",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envelope.NewEnvelope(actor, "benchmark-task", envelope.ProtocolMCP,
			"github.list_repos", "github.com/org", envelope.CapRead)
	}
}

func BenchmarkEnvelopeHash(b *testing.B) {
	env := benchEnvelope()
	env.PolicyDecision = envelope.DecisionAllow
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		env.Hash()
	}
}
