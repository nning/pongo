package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/ipv6"
)

type Net struct {
	config *Config
	ID     string
	peers  map[string]*net.Addr
}

func (n *Net) announce(iface string) {
	// log.Printf("announce %s on %s", n.ID, iface)

	daddr, err := net.ResolveUDPAddr("udp6", "[ff12::7179%"+iface+"]:7179")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp6", nil, daddr)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	conn.Write([]byte(n.ID))
}

func (n *Net) Announce() {
	for {
		for _, iface := range n.config.ListenInterfaces {
			n.announce(iface)
		}

		time.Sleep(time.Second)
	}
}

func (n *Net) listen(iface string) {
	addr := "[ff12::7179%" + iface + "]:7179"

	gaddr, err := net.ResolveUDPAddr("udp6", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenPacket("udp6", addr)
	if err != nil {
		log.Fatal(err)
	}

	pconn := ipv6.NewPacketConn(conn)

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, i := range ifaces {
		if i.Name == iface {
			if err := pconn.JoinGroup(&i, gaddr); err != nil {
				log.Fatal(err)
			}
		}
	}

	bs := make([]byte, 14)
	for {
		c, _, paddr, err := pconn.ReadFrom(bs)
		if err != nil {
			log.Fatal(err)
		}

		// log.Printf("recv %d bytes from %s", c, paddr)

		if c != 14 {
			continue
		}

		id := string(bs)
		if id == n.ID || n.peers[id] != nil {
			continue
		}

		n.peers[id] = &paddr
		log.Printf("peer %s joined", id)
	}
}

func (n *Net) Listen() {
	for _, iface := range n.config.ListenInterfaces {
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
