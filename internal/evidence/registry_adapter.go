package evidence

import (
	"encoding/json"
	"fmt"
)

// RegistryAdminAdapter exposes a ChainRegistry through the same admin interface
// the single-chain AdminAdapter uses, but resolves the right per-session chain
// for each call.
type RegistryAdminAdapter struct {
	reg *ChainRegistry
}

func NewRegistryAdminAdapter(reg *ChainRegistry) *RegistryAdminAdapter {
	return &RegistryAdminAdapter{reg: reg}
}

func (a *RegistryAdminAdapter) ExportSession(sessionID string) (interface{}, error) {
	chain, err := a.reg.get(sessionID)
	if err != nil {
		return nil, err
	}
	if chain == nil {
		return nil, fmt.Errorf("session %q not found", sessionID)
	}
	bundle, err := chain.Export()
	if err != nil {
		return nil, err
	}
	return json.RawMessage(bundle), nil
}

func (a *RegistryAdminAdapter) VerifySession(sessionID string) (interface{}, error) {
	chain, err := a.reg.get(sessionID)
	if err != nil {
		return nil, err
	}
	if chain == nil {
		return nil, fmt.Errorf("session %q not found", sessionID)
	}
	// Verify signatures when the registry is keyed, otherwise just the structure.
	if key := a.reg.Key(); len(key) > 0 {
		return VerifySignatures(chain.Records(), key), nil
	}
	return Verify(chain.Records()), nil
}

func (a *RegistryAdminAdapter) ListSessions() (interface{}, error) {
	chains, err := a.reg.all()
	if err != nil {
		return nil, err
	}
	manifests := make([]SessionManifest, 0, len(chains))
	for _, c := range chains {
		manifests = append(manifests, c.Manifest())
	}
	return manifests, nil
}

func (a *RegistryAdminAdapter) RenderReport(sessionID string) (string, error) {
	chain, err := a.reg.get(sessionID)
	if err != nil {
		return "", err
	}
	if chain == nil {
		return "", fmt.Errorf("session %q not found", sessionID)
	}
	return RenderMarkdownReport(chain)
}

func (a *RegistryAdminAdapter) RenderHTMLReport(sessionID string) (string, error) {
	chain, err := a.reg.get(sessionID)
	if err != nil {
		return "", err
	}
	if chain == nil {
		return "", fmt.Errorf("session %q not found", sessionID)
	}
	return RenderHTMLReport(chain)
}
