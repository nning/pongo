package main

import (
	"log"
	"net"

	"golang.org/x/net/ipv6"
)

var lastStateRecv []byte

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

		msg, err := decode[Announce](bs[:c])
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

		seq := bs[0]
		log.Println(seq)

		if seq == 0 {
			lastStateRecv = bs[1:c]

			msg, err := decode[State](bs[1:c])
			if err != nil {
				log.Println(err)
				continue
			}
			if msg.ID != id {
				continue
			}

			log.Printf("recv full %v\n", msg)

			n.LastState = msg.Game
		} else {
			d, err := decode[Diff](bs[1:c])
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("recv diff len %v\n", len(*d))

			if len(lastStateRecv) == 0 {
				continue
			}

			mbs := patch(lastStateRecv, *d)
			if len(mbs) == 0 {
				log.Fatal("patch failed")
			}

			msg, err := decode[State](mbs)
			if err != nil {
				log.Println(err)
				continue
			}
			if msg.ID != id {
				continue
			}

			log.Printf("recv diff %v\n", msg)
		}
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
