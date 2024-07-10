package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/packet/bgp"
	"github.com/osrg/gobgp/v3/pkg/server"
	apb "google.golang.org/protobuf/types/known/anypb"
)

type DeviceConfig struct {
	DeviceID string `json:"device_id"`
	Config   struct {
		Interface string `json:"interface"`
		IPAddress string `json:"ip_address"`
		Netmask   string `json:"netmask"`
	} `json:"config"`
}

func main() {
	s := server.NewBgpServer()
	go s.Serve()

	// start a bgp server. ips are from docker networking
	if err := s.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: &api.Global{
			Asn:             65003,
			RouterId:        "172.20.0.2",
			ListenPort:      179,
			ListenAddresses: []string{"172.20.0.2"},
		},
	}); err != nil {
		log.Fatal(err)
	}

	n := &api.Peer{
		Conf: &api.PeerConf{
			NeighborAddress: "172.20.0.3",
			PeerAsn:         65002,
		},
		Transport: &api.Transport{
			LocalAddress: "172.20.0.2",
		},
	}

	if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
		Peer: n,
	}); err != nil {
		log.Fatal(err)
	}

	// Define the configuration data. In my crazy example I am
	// passing through just some configuration for a interface on
	configData := DeviceConfig{
		DeviceID: "router1",
		Config: struct {
			Interface string `json:"interface"`
			IPAddress string `json:"ip_address"`
			Netmask   string `json:"netmask"`
		}{
			Interface: "GigabitEthernet0/1",
			IPAddress: "192.168.1.1",
			Netmask:   "255.255.255.0",
		},
	}

	// Serialize the configuration data to JSON
	configBytes, err := json.Marshal(configData)
	if err != nil {
		log.Fatalf("Failed to serialize configuration data: %s\n", err)
	}

	// Create a custom BGP attribute for the configuration data
	nlri, _ := apb.New(&api.IPAddressPrefix{
		PrefixLen: 24,
		Prefix:    "192.168.1.0",
	})

	// Add ORIGIN attribute
	origin, _ := apb.New(&api.OriginAttribute{
		Origin: 0, // 0 means IGP
	})

	// Add NEXT_HOP attribute
	nextHop, _ := apb.New(&api.NextHopAttribute{
		NextHop: "172.20.0.2",
	})

	a1, _ := apb.New(&api.UnknownAttribute{
		Flags: uint32(bgp.BGP_ATTR_FLAG_OPTIONAL | bgp.BGP_ATTR_FLAG_TRANSITIVE),
		Type:  99,
		Value: configBytes,
	})

	// Include origin, nextHop, and custom attribute
	attribs := []*apb.Any{origin, nextHop, a1}

	log.Printf("Sending path update with attributes: %+v\n", attribs)

	_, err = s.AddPath(context.Background(), &api.AddPathRequest{
		Path: &api.Path{
			Family: &api.Family{Afi: api.Family_AFI_IP, Safi: api.Family_SAFI_UNICAST},
			Nlri:   nlri,
			Pattrs: attribs,
		},
	})
	if err != nil {
		log.Fatalf("Failed to add path: %s\n", err)
	}

	fmt.Println("BGP speaker is running...")
	select {}
}
