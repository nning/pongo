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
	"time"

	"golang.org/x/net/ipv6"
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

func (msg State) String() string {
	return fmt.Sprintf("{%v, Ball:{%v, %v}, Paddle1:{%v}, Paddle2:{%v}}", msg.ID, msg.Game.Ball.X, msg.Game.Ball.Y, msg.Game.Paddle1.Y, msg.Game.Paddle2.Y)
}

func (n *Net) sendState(game *Game) {
	if n.peer == nil {
		return
	}

	// Copy Addr, use announce port as remote target port
	var addr net.UDPAddr = *n.peer
	addr.Port = n.announcePort

	conn, err := net.DialUDP("udp6", nil, &addr)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err = enc.Encode(State{n.id, game})
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("send state %v (size %v)", n.peer, c)
}

func (n *Net) SendState(game *Game) {
	for {
		n.sendState(game)
		time.Sleep(time.Second / 60)
	}
}

func (n *Net) announce(iface net.Interface) {
	daddr, err := net.ResolveUDPAddr("udp6", "[ff12::7179%"+iface.Name+"]:7179")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp6", nil, daddr)
	if err != nil {
		return
	}
	defer conn.Close()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	err = enc.Encode(Announce{n.id, n.announcePort})
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("announce %v on %v (port %v)\n", n.id, iface.Name, n.announcePort)
}

func (n *Net) Announce() {
	for {
		ifaces, err := net.Interfaces()
		if err != nil {
			log.Fatal(err)
		}

		for _, iface := range ifaces {
			n.announce(iface)
		}

		time.Sleep(time.Second / 3)
	}
}

func (n *Net) listenAnnounce(iface net.Interface) {
	addr := "[ff12::7179%" + iface.Name + "]:7179"

	gaddr, err := net.ResolveUDPAddr("udp6", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("udp6", addr)
	if err != nil {
		log.Fatal(err)
	}

	pconn := ipv6.NewPacketConn(conn)
	if err := pconn.JoinGroup(&iface, gaddr); err != nil {
		log.Fatal(err)
	}

	bs := make([]byte, 256)
	for {
		if !n.announceEnabled {
			break
		}

		c, _, paddr, err := pconn.ReadFrom(bs)
		if err != nil {
			log.Fatal(err)
		}

		buf := bytes.NewBuffer(bs[:c])
		dec := gob.NewDecoder(buf)

		var msg Announce
		err = dec.Decode(&msg)
		if err != nil {
			log.Println(err)
			continue
		}

		n.peerMutex.Lock()
		if msg.ID == n.id || n.peer != nil {
			n.peerMutex.Unlock()
			continue
		}

		if n.announcePort == msg.Port {
			log.Fatalf("port %v collision", msg.Port)
		} else if n.announcePort < msg.Port {
			game.mode = Left
		} else {
			game.mode = Right
		}

		n.peer = paddr.(*net.UDPAddr)
		n.peer.Port = msg.Port

		n.peerMutex.Unlock()

		log.Printf("peer %s joined: %s (we are %v)", msg.ID, n.peer, game.mode)
		n.announceEnabled = false // Concurrence safe?

		go n.listenState(msg.ID)
	}
}

func (n *Net) listenState(id string) {
	log.Printf("listen for state from %v\n", id)

	var addr net.UDPAddr = *n.peer
	addr.IP = net.IPv6zero

	conn, err := net.ListenPacket("udp6", addr.String())
	if err != nil {
		log.Fatal(err)
	}

	bs := make([]byte, 512)
	for {
		c, _, err := conn.ReadFrom(bs)
		if err != nil {
			log.Fatal(err)
		}

		buf := bytes.NewBuffer(bs[:c])
		dec := gob.NewDecoder(buf)

		var msg State
		err = dec.Decode(&msg)
		if err != nil {
			log.Println(err)
			continue
		}
		if msg.ID != id {
			continue
		}

		// log.Printf("recv from %v: %v\n", id, msg)

		n.LastState = msg.Game
	}
}

func (n *Net) Listen() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifaces {
		go n.listenAnnounce(iface)
	}
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
