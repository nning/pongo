package main

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"net"
	"sync"
)

type Net struct {
	config          *Config
	id              string
	peer            *net.UDPAddr
	peerMutex       sync.Mutex
	announceEnabled bool
	announcePort    int
	LastState       *Game
}

type Announce struct {
	ID   string
	Port int
}

type State struct {
	ID   string
	Game *Game
}

type Diff map[int]byte

func (msg State) String() string {
	return fmt.Sprintf("{%v, Ball: {%v, %v}, Paddle1: {%v}, Paddle2: {%v}}", msg.ID, msg.Game.Ball.X, msg.Game.Ball.Y, msg.Game.Paddle1.Y, msg.Game.Paddle2.Y)
}

// diff returns the changed bytes in inc compared to base as map of byte index
// to changed byte value
func diff(base, inc []byte) *Diff {
	i := 0
	m := make(Diff)

	for ; i < len(base) && i < len(inc); i++ {
		if base[i] != inc[i] {
			m[i] = inc[i]
		}
	}

	if i < len(inc) {
		for ; i < len(inc); i++ {
			m[i] = inc[i]
		}
	}

	return &m
}

func patch(base []byte, diff *Diff) []byte {
	bs := make([]byte, len(base))
	copy(bs, base)

	for i, v := range *diff {
		if i < len(bs) {
			bs[i] = v
		} else {
			bs = append(bs, (*diff)[i])
		}
	}

	return bs
}

func decode[T any](bs []byte) (*T, error) {
	buf := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(buf)

	var t T
	err := dec.Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func encode[T any](t *T) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(t)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getRandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x", b[0:2], b[2:4], b[4:6])
}

func NewNet(config *Config) *Net {
	rp, err := rand.Int(rand.Reader, big.NewInt(65534-1024))
	if err != nil {
		log.Fatal(err)
	}

	p := int(rp.Int64()) + 1024
	if config.ListenPort != 0 {
		p = config.ListenPort
	}

	return &Net{
		config:          config,
		id:              getRandomID(),
		peerMutex:       sync.Mutex{},
		announceEnabled: true,
		announcePort:    p,
	}
}
