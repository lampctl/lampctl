package sequencer

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lampctl/lampctl/registry"
	"gitlab.com/gomidi/midi/v2/smf"
)

type sequencerEntry struct {
	Provider *registry.Provider
	Changes  []*registry.Change
}

type sequencerGroup struct {
	Offset  time.Duration
	Entries []*sequencerEntry
}

type sequencerSequence struct {
	Groups     []*sequencerGroup
	GroupIndex int
}

type mappingNote struct {
	ProviderID string `json:"provider_id"`
	GroupID    string `json:"group_id"`
	LampID     string `json:"lamp_id"`
}

type mappingMap map[string]*mappingNote

func (s *Sequencer) loadMap(mappingFilename string) (mappingMap, error) {
	f, err := os.Open(mappingFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	m := make(mappingMap)
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Sequencer) load(midiFilename, mappingFilename string) error {
	f, err := smf.ReadFile(midiFilename)
	if err != nil {
		return err
	}
	m, err := s.loadMap(mappingFilename)
	if err != nil {
		return err
	}

	_ = m

	return nil
}
