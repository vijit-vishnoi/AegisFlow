package evidence

import (
	"database/sql"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/saivedant169/AegisFlow/internal/envelope"
)

// Recorder records an action into an evidence chain. Both a single SessionChain
// and a ChainRegistry satisfy it, so callers (e.g. the MCP gateway) don't care
// whether records go to one chain or are split per session.
type Recorder interface {
	Record(env *envelope.ActionEnvelope) (*Record, error)
}

// ChainRegistry keeps one evidence chain per active session. Idle chains are
// evicted from memory. Persistent mode reloads an evicted chain from SQLite.
type ChainRegistry struct {
	mu      sync.Mutex
	chains  map[string]*registryEntry
	key     []byte
	maxIdle time.Duration
	stop    chan struct{}
	store   *sqliteRecordStore
}

type registryEntry struct {
	chain    *SessionChain
	lastUsed time.Time
}

const defaultSessionFallback = "default"

// NewChainRegistry creates a registry whose chains are signed with key (may be
// nil for unsigned). It starts a janitor that evicts sessions idle longer than
// maxIdle; pass 0 for a sensible default.
func NewChainRegistry(key []byte) *ChainRegistry {
	r := &ChainRegistry{
		chains:  make(map[string]*registryEntry),
		key:     append([]byte(nil), key...),
		maxIdle: time.Hour,
		stop:    make(chan struct{}),
	}
	go r.janitor()
	return r
}

// NewPersistentChainRegistry restores signed evidence from SQLite and rejects
// startup when stored records fail integrity or signature verification.
func NewPersistentChainRegistry(key []byte, db *sql.DB) (*ChainRegistry, error) {
	store, err := newSQLiteRecordStore(db, key)
	if err != nil {
		return nil, err
	}
	r := &ChainRegistry{
		chains:  make(map[string]*registryEntry),
		key:     append([]byte(nil), key...),
		maxIdle: time.Hour,
		stop:    make(chan struct{}),
		store:   store,
	}
	ids, err := store.ListSessionIDs()
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		records, err := store.LoadSession(id)
		if err != nil {
			return nil, fmt.Errorf("restore evidence session %q: %w", id, err)
		}
		_, err = restoreSessionChain(id, r.key, records, nil)
		if err != nil {
			return nil, fmt.Errorf("verify evidence session %q: %w", id, err)
		}
	}
	go r.janitor()
	return r, nil
}

func (r *ChainRegistry) janitor() {
	t := time.NewTicker(r.maxIdle / 4)
	defer t.Stop()
	for {
		select {
		case <-r.stop:
			return
		case now := <-t.C:
			r.mu.Lock()
			for id, e := range r.chains {
				if now.Sub(e.lastUsed) > r.maxIdle {
					delete(r.chains, id)
				}
			}
			r.mu.Unlock()
		}
	}
}

// chainFor returns the chain for a session, creating a signed one on first use.
func (r *ChainRegistry) chainFor(sessionID string) (*SessionChain, error) {
	if sessionID == "" {
		sessionID = defaultSessionFallback
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	e, ok := r.chains[sessionID]
	if !ok {
		if r.store != nil {
			records, err := r.store.LoadSession(sessionID)
			if err != nil {
				return nil, fmt.Errorf("load evidence session %q: %w", sessionID, err)
			}
			if len(records) > 0 {
				chain, err := restoreSessionChain(sessionID, r.key, records, func(record Record) error {
					return r.store.Append(sessionID, record)
				})
				if err != nil {
					return nil, fmt.Errorf("verify evidence session %q: %w", sessionID, err)
				}
				e = &registryEntry{chain: chain}
			}
		}
		if e == nil {
			chain := NewSignedSessionChain(sessionID, r.key)
			if r.store != nil {
				chain.append = func(record Record) error {
					return r.store.Append(sessionID, record)
				}
			}
			e = &registryEntry{chain: chain}
		}
		r.chains[sessionID] = e
	}
	e.lastUsed = time.Now()
	return e.chain, nil
}

// Record routes the action to its session's chain.
func (r *ChainRegistry) Record(env *envelope.ActionEnvelope) (*Record, error) {
	if env == nil {
		return nil, fmt.Errorf("evidence envelope is required")
	}
	chain, err := r.chainFor(env.Actor.SessionID)
	if err != nil {
		return nil, err
	}
	return chain.Record(env)
}

// get returns the chain for a session without creating one.
func (r *ChainRegistry) get(sessionID string) (*SessionChain, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if e, ok := r.chains[sessionID]; ok {
		e.lastUsed = time.Now()
		return e.chain, nil
	}
	if r.store == nil {
		return nil, nil
	}
	records, err := r.store.LoadSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("load evidence session %q: %w", sessionID, err)
	}
	if len(records) == 0 {
		return nil, nil
	}
	chain, err := restoreSessionChain(sessionID, r.key, records, func(record Record) error {
		return r.store.Append(sessionID, record)
	})
	if err != nil {
		return nil, fmt.Errorf("verify evidence session %q: %w", sessionID, err)
	}
	r.chains[sessionID] = &registryEntry{chain: chain, lastUsed: time.Now()}
	return chain, nil
}

// chains snapshot for listing.
func (r *ChainRegistry) all() ([]*SessionChain, error) {
	if r.store != nil {
		ids, err := r.store.ListSessionIDs()
		if err != nil {
			return nil, err
		}
		for _, id := range ids {
			if _, err := r.get(id); err != nil {
				return nil, err
			}
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*SessionChain, 0, len(r.chains))
	for _, e := range r.chains {
		out = append(out, e.chain)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].SessionID() < out[j].SessionID()
	})
	return out, nil
}

// Key returns the signing key (used by callers that need to verify signatures).
func (r *ChainRegistry) Key() []byte { return append([]byte(nil), r.key...) }

// Close stops the eviction janitor.
func (r *ChainRegistry) Close() {
	select {
	case <-r.stop:
	default:
		close(r.stop)
	}
}
