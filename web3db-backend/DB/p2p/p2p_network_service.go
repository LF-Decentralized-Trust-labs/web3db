package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	peerstore "github.com/libp2p/go-libp2p/core/peer"
	routing "github.com/libp2p/go-libp2p/core/routing"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	connmgr "github.com/libp2p/go-libp2p/p2p/net/connmgr"
	noise "github.com/libp2p/go-libp2p/p2p/security/noise"
	tls "github.com/libp2p/go-libp2p/p2p/security/tls"
	discovery "github.com/libp2p/go-libp2p/discovery"
)

func setupDiscovery(ctx context.Context, h host.Host, dht *dht.IpfsDHT) {
	// Setup mDNS discovery
	mdnsService, err := mdns.NewMdnsService(ctx, h, time.Hour, "web3db-service")
	if err != nil {
		log.Fatalf("failed to start mDNS: %v", err)
	}
	defer mdnsService.Close()

	// Register with the local network
	mdns.DiscoverLocalPeers(ctx, "web3db-service", h)

	// Use DHT for global peer discovery
	routingDiscovery := discovery.NewRoutingDiscovery(dht)
	discovery.Advertise(ctx, routingDiscovery, "web3db-service")

	peerChan, err := routingDiscovery.FindPeers(ctx, "web3db-service")
	if err != nil {
		log.Fatalf("failed to find peers: %v", err)
	}

	// Handle discovered peers
	go func() {
		for p := range peerChan {
			if p.ID == h.ID() {
				continue
			}
			log.Printf("Discovered new peer %s\n", p.ID.Pretty())
			if _, err := h.Network().DialPeer(ctx, p.ID); err != nil {
				log.Printf("Error connecting to peer %s: %v", p.ID.Pretty(), err)
			}
		}
	}()
}

func main() {
	ctx := context.Background()

	// Generate a new identity key pair
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, rand.Reader)
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}

	connManager, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		log.Fatalf("failed to create the connection manager: %v", err)
	}

	// Create a libp2p host
	h, err := libp2p.New(ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"),
		libp2p.Security(noise.ID, noise.New),
		libp2p.Security(tls.ID, tls.New),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connManager),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.New(ctx, h, dht.Mode(dht.ModeServer))
		}),
	)
	if err != nil {
		log.Fatalf("failed to create host: %v", err)
	}
	defer h.Close()

	// Setup DHT for discovery
	dht, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	if err != nil {
		log.Fatalf("failed to create DHT: %v", err)
	}

	setupDiscovery(ctx, h, dht)

	log.Printf("Your Host ID is: %s", h.ID())
	select {} // hang forever
}
