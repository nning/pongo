package main

import (
	"testing"

	"github.com/fxamacker/cbor/v2"
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

func TestDiff(t *testing.T) {
	bs1 := []byte{}
	d1 := diff(bs1, bs1)
	assert.Equal(t, 0, len(*d1))
	assert.Equal(t, &Diff{}, d1)

	bs1 = []byte{1, 2, 3}
	bs2 := []byte{}
	d1 = diff(bs1, bs2)
	assert.Equal(t, 0, len(*d1))
	assert.Equal(t, &Diff{}, d1)

	bs1 = []byte{1, 2, 3}
	bs2 = []byte{1, 2, 4}
	d1 = diff(bs1, bs2)
	assert.Equal(t, 1, len(*d1))
	assert.Equal(t, &Diff{2: 4}, d1)

	bs1 = []byte{}
	bs2 = []byte{1, 2, 3}
	d1 = diff(bs1, bs2)
	assert.Equal(t, 3, len(*d1))
	assert.Equal(t, &Diff{0: 1, 1: 2, 2: 3}, d1)

	bs1 = []byte{1, 2, 3}
	bs2 = []byte{1, 2}
	d1 = diff(bs1, bs2)
	assert.Equal(t, 1, len(*d1))
	assert.Equal(t, &Diff{-2: 0}, d1)

	bs1 = []byte{1, 2, 3, 4, 5}
	bs2 = []byte{1, 2}
	d1 = diff(bs1, bs2)
	assert.Equal(t, 1, len(*d1))
	assert.Equal(t, &Diff{-2: 0}, d1)
}

func TestPatch(t *testing.T) {
	bs1 := []byte{}
	bs2 := patch(bs1, &Diff{0: 1, 1: 2, 2: 3})
	assert.Equal(t, []byte{1, 2, 3}, bs2)

	bs1 = []byte{1, 2, 3}
	bs2 = patch(bs1, &Diff{2: 4})
	assert.Equal(t, []byte{1, 2, 4}, bs2)

	bs1 = []byte{1, 2, 3}
	bs2 = patch(bs1, &Diff{0: 2, 2: 2})
	assert.Equal(t, []byte{2, 2, 2}, bs2)

	bs1 = []byte{1, 2, 3}
	bs2 = patch(bs1, &Diff{-2: 0})
	assert.Equal(t, []byte{1, 2}, bs2)

	bs1 = []byte{1, 2, 3, 4, 5}
	bs2 = patch(bs1, &Diff{-2: 0})
	assert.Equal(t, []byte{1, 2}, bs2)
}

func TestDecode(t *testing.T) {
	buf, err := cbor.Marshal(state)
	assert.Nil(t, err)

	msg, err := decode[State](buf)
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

	d1 := diff(bs1, bs2)
	assert.Equal(t, 2, len(*d1))
	assert.Equal(t, &Diff{34: 0x40, 35: 0x33}, d1)

	bs3, err := encode(d1)
	assert.Nil(t, err)

	d2, err := decode[Diff](bs3)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(*d2))
	assert.Equal(t, &Diff{34: 0x40, 35: 0x33}, d2)

	bs4 := patch(bs1, d2)
	assert.NotEqual(t, 0, len(bs4))

	s, err = decode[State](bs4)
	assert.Nil(t, err)
	assert.Equal(t, float64(19), s.Game.Ball.X)
	assert.Equal(t, float64(2), s.Game.Ball.Y)
}
