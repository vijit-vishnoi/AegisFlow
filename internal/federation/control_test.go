package federation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/saivedant169/AegisFlow/internal/config"
)

func TestConfigHandlerStripsSecrets(t *testing.T) {
	cfg := &config.Config{
		Providers: []config.ProviderConfig{{
			Name:    "provider",
			APIKeys: []config.ProviderAPIKey{{Key: "provider-secret", KeyEnv: "PROVIDER_KEY"}},
			Config: map[string]string{
				"api_token": "provider-config-secret",
				"region":    "us-test-1",
			},
		}},
		Tenants: []config.TenantConfig{{
			ID:      "tenant-1",
			APIKeys: []config.APIKeyEntry{{Key: "tenant-secret", Role: "admin"}},
		}},
		Federation: config.FederationConfig{
			DataPlanes: []config.DataPlaneConfig{{Name: "plane-1", Token: "federation-secret"}},
		},
		Credentials: config.CredentialConfig{Providers: []config.CredentialProviderCfg{{
			Name:          "broker",
			Token:         "broker-secret",
			GitHubKeyPath: "/secret/key.pem",
			VaultToken:    "vault-secret",
			AWSExternalID: "aws-secret",
		}}},
	}
	cfg.RateLimit.Redis.Password = "rate-limit-secret"
	cfg.Cache.Redis.Password = "cache-secret"
	cfg.Cache.Semantic.APIKey = "semantic-secret"
	cfg.Webhook.Secret = "webhook-secret"
	cfg.Database.ConnString = "postgres://user:database-secret@db/service"
	cfg.Admin.Token = "admin-secret"
	cfg.Eval.Webhook.URL = "https://eval.example/secret-path"
	cfg.ApprovalIntegrations.GitHub.Token = "github-secret"
	cfg.ApprovalIntegrations.Slack.WebhookURL = "https://hooks.example/slack-secret"
	cfg.Capability.SigningKey = "capability-secret"
	cfg.SupplyChain.SigningKey = "supply-chain-secret"

	cp := NewControlPlane(cfg)
	req := httptest.NewRequest(http.MethodGet, "/admin/v1/federation/config", nil)
	req.Header.Set("Authorization", "Bearer federation-secret")
	recorder := httptest.NewRecorder()

	cp.ConfigHandler(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, secret := range []string{
		"provider-secret",
		"provider-config-secret",
		"tenant-secret",
		"federation-secret",
		"broker-secret",
		"/secret/key.pem",
		"vault-secret",
		"aws-secret",
		"rate-limit-secret",
		"cache-secret",
		"semantic-secret",
		"webhook-secret",
		"database-secret",
		"admin-secret",
		"secret-path",
		"github-secret",
		"slack-secret",
		"capability-secret",
		"supply-chain-secret",
	} {
		if strings.Contains(body, secret) {
			t.Errorf("response leaked %q", secret)
		}
	}
	if !strings.Contains(body, "us-test-1") || !strings.Contains(body, "PROVIDER_KEY") {
		t.Fatalf("response removed non-secret provider config: %s", body)
	}
}

func TestConfigHandlerRejectsUnknownToken(t *testing.T) {
	cp := NewControlPlane(&config.Config{Federation: config.FederationConfig{
		DataPlanes: []config.DataPlaneConfig{{Name: "plane-1", Token: "expected"}},
	}})
	req := httptest.NewRequest(http.MethodGet, "/admin/v1/federation/config", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	recorder := httptest.NewRecorder()

	cp.ConfigHandler(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", recorder.Code)
	}
}
