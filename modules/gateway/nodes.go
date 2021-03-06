package gateway

import (
	"errors"
	"math/rand"

	"github.com/NebulousLabs/Sia/encoding"
	"github.com/NebulousLabs/Sia/modules"
)

const (
	maxSharedNodes = 10
	maxAddrLength  = 100
	minPeers       = 3
)

// addNode adds an address to the set of nodes on the network.
func (g *Gateway) addNode(addr modules.NetAddress) error {
	if _, exists := g.nodes[addr]; exists {
		return errors.New("node already added")
	}
	g.nodes[addr] = struct{}{}
	g.log.Println("INFO: added node", addr)
	return nil
}

func (g *Gateway) removeNode(addr modules.NetAddress) error {
	if _, exists := g.nodes[addr]; !exists {
		return errors.New("no record of that node")
	}
	delete(g.nodes, addr)
	g.log.Println("INFO: removed node", addr)
	return nil
}

func (g *Gateway) randomNode() (modules.NetAddress, error) {
	if len(g.nodes) > 0 {
		r := rand.Intn(len(g.nodes))
		for node := range g.nodes {
			if r == 0 {
				return node, nil
			}
			r--
		}
	}

	return "", errNoPeers
}

// shareNodes is an RPC that returns up to 10 randomly selected nodes.
func (g *Gateway) shareNodes(conn modules.PeerConn) error {
	id := g.mu.RLock()
	var nodes []modules.NetAddress
	for node := range g.nodes {
		if len(nodes) == maxSharedNodes {
			break
		}
		nodes = append(nodes, node)
	}
	g.mu.RUnlock(id)
	return encoding.WriteObject(conn, nodes)
}

// relayNode adds a node to the Gateway's node list and relays it to each of
// the Gateway's peers. If the node is already in the node list, it is not
// relayed.
func (g *Gateway) relayNode(conn modules.PeerConn) error {
	// read address
	var addr modules.NetAddress
	if err := encoding.ReadObject(conn, &addr, maxAddrLength); err != nil {
		return err
	}
	// add node
	id := g.mu.Lock()
	err := g.addNode(addr)
	g.mu.Unlock(id)
	// relay
	if err == nil {
		go g.Broadcast("RelayNode", addr)
	}
	return nil
}
