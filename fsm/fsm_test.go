package fsm

import (
	"testing"

	"github.com/dc0d/club/errors"
	"github.com/stretchr/testify/assert"
)

var (
	errNegative = errors.Errorf("NEGATIVE")
)

type sample struct {
	state int
}

func (s *sample) Start() State {
	return s.dispatcher
}

func (s *sample) odd() (State, error) {
	s.state++
	return s.dispatcher, nil
}

func (s *sample) even() (State, error) {
	s.state++
	return s.dispatcher, nil
}

func (s *sample) dispatcher() (State, error) {
	if s.state > 8 {
		return nil, nil
	}
	if s.state < 0 {
		return nil, errNegative
	}
	if s.state%2 == 0 {
		return s.even, nil
	}
	return s.odd, nil
}

func Test01(t *testing.T) {
	assert := assert.New(t)

	fsm := &sample{}
	Activate(fsm.Start())

	assert.Equal(9, fsm.state)
}
