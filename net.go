package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv6"
)

type Net struct {
	config *Config
	ID     string
	peers  map[string]*net.Addr
}

var peersMutex = &sync.Mutex{}

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

	conn.Write([]byte(n.ID))
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

		time.Sleep(time.Second)
	}
}

func (n *Net) listen(iface net.Interface) {
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

	bs := make([]byte, 14)
	for {
		c, _, paddr, err := pconn.ReadFrom(bs)
		if err != nil {
			log.Fatal(err)
		}

		if c != 14 {
			continue
		}

		id := string(bs)
		peersMutex.Lock()
		if id == n.ID || n.peers[id] != nil {
			peersMutex.Unlock()
			continue
		}

		n.peers[id] = &paddr
		peersMutex.Unlock()

		log.Printf("peer %s joined: %s", id, paddr.String())
	}
}

func (n *Net) Listen() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range ifaces {
		go n.listen(iface)
	}
}

func getRandomID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x", b[0:2], b[2:4], b[4:6])
}

func NewNet(config *Config) *Net {
	return &Net{
		config: config,
		ID:     getRandomID(),
		peers:  make(map[string]*net.Addr),
	}
}
