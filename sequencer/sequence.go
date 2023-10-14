package sequencer

import (
	"encoding/json"
	"os"
	"time"

	"github.com/lampctl/lampctl/registry"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/smf"
)

type sequencerRawEvent struct {
	Offset time.Duration
	Note   int
	NoteOn bool
}

type sequencerEvent struct {
	Provider *registry.Provider
	Changes  []*registry.Change
}

type sequencerGroup struct {
	Offset time.Duration
	Events []*sequencerEvent
}

type sequencerSequence struct {
	Groups     []*sequencerGroup
	GroupIndex int
}

func (s *Sequencer) loadRawEvents(midiFilename string) ([]*sequencerRawEvent, error) {
	f, err := smf.ReadFile(midiFilename)
	if err != nil {
		return nil, err
	}
	events := []*sequencerRawEvent{}
	for _, track := range f.Tracks {
		var t int64
		for _, e := range track {
			t += int64(e.Delta)
			if !e.Message.IsOneOf(midi.NoteOnMsg, midi.NoteOffMsg) {
				continue
			}
			var (
				absOffset              = time.Duration(f.TimeAt(t) * 1000)
				channel, key, velocity uint8
			)
			switch e.Message.Type() {
			case midi.NoteOnMsg:
				e.Message.GetNoteOn(&channel, &key, &velocity)
			case midi.NoteOffMsg:
				e.Message.GetNoteOff(&channel, &key, &velocity)
			}
			events = append(events, &sequencerRawEvent{
				Offset: absOffset,
				Note:   int(key),
				NoteOn: e.Message.Is(midi.NoteOnMsg),
			})
		}
	}
	return events, nil
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

	// Read the raw MIDI events
	e, err := s.loadRawEvents(midiFilename)
	if err != nil {
		return err
	}

	// Read the mapping file
	m, err := s.loadMap(mappingFilename)
	if err != nil {
		return err
	}

	// TODO: create an ordered sequence from the events and mapping
	_ = e
	_ = m

	return nil
}
