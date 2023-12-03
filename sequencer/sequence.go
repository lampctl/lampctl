package sequencer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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
	Provider registry.Provider
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

type changeMap map[registry.Provider][]*registry.Change

func (s *Sequencer) load(midiFilename, mappingFilename string) error {

	// Read the raw MIDI events
	events, err := s.loadRawEvents(midiFilename)
	if err != nil {
		return err
	}

	// Read the mapping file
	mapping, err := s.loadMap(mappingFilename)
	if err != nil {
		return err
	}

	// Create a map of provider IDs to actual Provider instances
	providerMap := map[string]registry.Provider{}
	for _, m := range mapping {
		if _, ok := providerMap[m.ProviderID]; !ok {
			p, err := s.registry.GetProvider(m.ProviderID)
			if err != nil {
				return err
			}
			providerMap[m.ProviderID] = p
		}
	}

	// Group the events by their offset and then provider
	var (
		sequence          = &sequencerSequence{}
		currentOffset     time.Duration
		changesByProvider changeMap
	)
	for _, e := range events {

		// If this is the first event or a new offset...
		if changesByProvider == nil || e.Offset != currentOffset {

			// Create a sequencerGroup for the events
			if changesByProvider != nil {
				g := &sequencerGroup{
					Offset: currentOffset,
				}
				for p, changeList := range changesByProvider {
					g.Events = append(g.Events, &sequencerEvent{
						Provider: p,
						Changes:  changeList,
					})
				}
				sequence.Groups = append(sequence.Groups, g)
			}

			// Reset the current offset and change map
			currentOffset = e.Offset
			changesByProvider = changeMap{}
		}

		// Find the mapping for the note
		m, ok := mapping[strconv.Itoa(e.Note)]
		if !ok {
			return fmt.Errorf("note %d has no mapping", e.Note)
		}

		// Find the provider in the map
		p, ok := providerMap[m.ProviderID]
		if !ok {
			return fmt.Errorf("provider %s does not exist", m.ProviderID)
		}

		// Add the events to the map
		changesByProvider[p] = append(
			changesByProvider[p],
			&registry.Change{
				GroupID: m.GroupID,
				LampID:  m.LampID,
				State:   e.NoteOn,
			},
		)
	}

	// Assign the sequence
	s.sequence = sequence

	return nil
}
