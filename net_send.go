package main

import (
	"log"
	"net"
	"time"
)

var lastStateSend []byte

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

	bs, err := encode(&Announce{n.id, n.announcePort})
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(bs)
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

func (n *Net) sendState(game *Game, seq int) {
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

	bs, err := encode(&State{n.id, game})
	if err != nil {
		log.Fatal(err)
	}

	if lastStateSend != nil && seq > 0 {
		d := diff(lastStateSend, bs)
		bs, err = encode(&d)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = conn.Write(append([]byte{byte(seq)}, bs...))
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("send state %v (size %v, diff %v)", n.peer, c, len(d1))

	if seq == 0 {
		lastStateSend = bs
	}
}

func (n *Net) SendState(game *Game) {
	s := 0

	for {
		n.sendState(game, s)

		time.Sleep(time.Second / 60)

		s += 1
		s %= 30 // send key state 2 times a second
	}
}
