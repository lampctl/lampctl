package sequencer

import (
	"github.com/lampctl/lampctl/registry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	commandLoad = iota
	commandPlay
	commandStop
)

type sequencerCmd struct {
	Command int
	Params  any
}

type sequencerCmdLoadParams struct {
	AudioFilename   string
	MidiFilename    string
	MappingFilename string
}

// Sequencer provides a means of playing a sequence (possibly loaded from disk)
// in realtime. A mapping file must also be provided to map MIDI
type Sequencer struct {
	logger     zerolog.Logger
	registry   *registry.Registry
	sequence   *sequencerSequence
	cmdChan    chan *sequencerCmd
	retChan    chan error
	closeChan  chan any
	closedChan chan any
}

func (s *Sequencer) run() {
	defer close(s.closedChan)
	defer s.logger.Info().Msg("sequencer stopped")
	s.logger.Info().Msg("sequencer started")
	var (
	//...
	)
	for {
		select {
		case c := <-s.cmdChan:
			switch c.Command {
			case commandLoad:
				p := c.Params.(*sequencerCmdLoadParams)
				s.retChan <- s.load(p.MidiFilename, p.MappingFilename)
			case commandPlay:
				break
			case commandStop:
				break
			}
		case <-s.closeChan:
			return
		}
	}
}

// New creates (but does not initialize) a new sequencer entry.
func New(cfg *Config) *Sequencer {
	s := &Sequencer{
		logger:     log.With().Str("package", "sequencer").Logger(),
		registry:   cfg.Registry,
		cmdChan:    make(chan *sequencerCmd),
		retChan:    make(chan error),
		closeChan:  make(chan any),
		closedChan: make(chan any),
	}
	go s.run()
	return s
}

// Load attempts to load the specified audio, MIDI, and mapping files.
func (s *Sequencer) Load(
	audioFilename, midiFilename, mappingFilename string,
) error {
	s.cmdChan <- &sequencerCmd{
		Command: commandLoad,
		Params: &sequencerCmdLoadParams{
			AudioFilename:   audioFilename,
			MidiFilename:    midiFilename,
			MappingFilename: mappingFilename,
		},
	}
	return <-s.retChan
}

// Play begins the loaded sequence.
func (s *Sequencer) Play() {
	s.cmdChan <- &sequencerCmd{
		Command: commandPlay,
	}
}

// Stop ends playback.
func (s *Sequencer) Stop() {
	s.cmdChan <- &sequencerCmd{
		Command: commandStop,
	}
}

// Close stops (if required) and shuts down the sequencer.
func (s *Sequencer) Close() {
	close(s.closeChan)
	<-s.closedChan
}
