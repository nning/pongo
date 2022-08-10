package main

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
)

var id = getRandomID()
var state = &State{
	ID: id,
	Game: &Game{
		Ball: &Ball{
			Rect{
				X: 1,
				Y: 2,
				W: 3,
				H: 4,
			},
			Vector{
				VX: 5,
				VY: 6,
			},
			7,
			8,
		},
		Paddle1: &Paddle{
			Rect{
				X: 9,
				Y: 10,
				W: 11,
				H: 12,
			},
			13,
		},
		Paddle2: &Paddle{
			Rect{
				X: 14,
				Y: 15,
				W: 16,
				H: 17,
			},
			18,
		},
	},
}

func TestDecode(t *testing.T) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(state)

	msg, err := decode[State](buf.Bytes())
	assert.Nil(t, err)
	assert.Equal(t, state, msg)
}

func TestEncode(t *testing.T) {
	bs, err := encode(state)
	assert.Nil(t, err)

	msg, err := decode[State](bs)
	assert.Nil(t, err)
	assert.Equal(t, state, msg)
}

func TestEncodeDiffPatchDecode(t *testing.T) {
	bs1, err := encode(state)
	assert.Nil(t, err)

	state.Game.Ball.X = 19
	bs2, err := encode(state)
	assert.Nil(t, err)

	assert.Equal(t, len(bs1), len(bs2))

	s, err := decode[State](bs2)
	assert.Nil(t, err)
	assert.Equal(t, float64(19), s.Game.Ball.X)

	d := diff(bs1, bs2)
	assert.Equal(t, 2, len(d))

	bs3 := patch(bs1, d)
	assert.NotEqual(t, 0, len(bs3))

	s, err = decode[State](bs3)
	assert.Nil(t, err)
	assert.Equal(t, float64(19), s.Game.Ball.X)
}
