package runner

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"

	"github.com/BurntSushi/toml"
)

type ObservedState struct {
	IP         net.IP
	ReverseDNS []string
}

func (state ObservedState) Equal(other ObservedState) bool {
	return state.IP.Equal(other.IP) && slices.Equal(state.ReverseDNS, other.ReverseDNS)
}

type StateStore interface {
	Load() (ObservedState, bool, error)
	Save(ObservedState) error
}

type memoryStateStore struct {
	state ObservedState
	known bool
}

func (store *memoryStateStore) Load() (ObservedState, bool, error) {
	return store.state, store.known, nil
}

func (store *memoryStateStore) Save(state ObservedState) error {
	store.state = state
	store.known = true
	return nil
}

type fileStateStore struct {
	path string
}

type stateFile struct {
	IP         string   `toml:"ip"`
	ReverseDNS []string `toml:"reverse_dns"`
}

func newStateStore(path string) StateStore {
	if path == "" {
		return &memoryStateStore{}
	}
	return fileStateStore{path: path}
}

func (store fileStateStore) Load() (ObservedState, bool, error) {
	data, err := os.ReadFile(store.path)
	if os.IsNotExist(err) {
		return ObservedState{}, false, nil
	}
	if err != nil {
		return ObservedState{}, false, fmt.Errorf("read state file: %w", err)
	}

	var encoded stateFile
	if err := toml.Unmarshal(data, &encoded); err != nil {
		return ObservedState{}, false, fmt.Errorf("decode state file: %w", err)
	}
	ip := net.ParseIP(encoded.IP)
	if ip == nil {
		return ObservedState{}, false, fmt.Errorf("state file contains malformed IP %q", encoded.IP)
	}
	return ObservedState{IP: ip, ReverseDNS: encoded.ReverseDNS}, true, nil
}

func (store fileStateStore) Save(state ObservedState) error {
	var buffer bytes.Buffer
	if err := toml.NewEncoder(&buffer).Encode(stateFile{
		IP:         state.IP.String(),
		ReverseDNS: state.ReverseDNS,
	}); err != nil {
		return fmt.Errorf("encode state file: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(store.path), 0o755); err != nil {
		return fmt.Errorf("create state file directory: %w", err)
	}
	if err := os.WriteFile(store.path, buffer.Bytes(), 0o600); err != nil {
		return fmt.Errorf("write state file: %w", err)
	}
	return nil
}
