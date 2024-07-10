package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	api "github.com/osrg/gobgp/v3/api"
	"github.com/osrg/gobgp/v3/pkg/server"
	"google.golang.org/protobuf/proto"
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

	if err := s.StartBgp(context.Background(), &api.StartBgpRequest{
		Global: &api.Global{
			Asn:             65002,
			RouterId:        "172.20.0.3",
			ListenPort:      -1,
			ListenAddresses: []string{"172.20.0.3"},
		},
	}); err != nil {
		log.Fatal(err)
	}

	n := &api.Peer{
		Conf: &api.PeerConf{
			NeighborAddress: "172.20.0.2",
			PeerAsn:         65003,
		},
		Transport: &api.Transport{
			LocalAddress: "172.20.0.3",
		},
	}

	if err := s.AddPeer(context.Background(), &api.AddPeerRequest{
		Peer: n,
	}); err != nil {
		log.Fatal(err)
	}

	if err := s.WatchEvent(context.Background(), &api.WatchEventRequest{
		Table: &api.WatchEventRequest_Table{
			Filters: []*api.WatchEventRequest_Table_Filter{
				{
					Type: api.WatchEventRequest_Table_Filter_BEST,
				},
			},
		},
	}, func(r *api.WatchEventResponse) {
		if table := r.GetTable(); table != nil {
			for _, path := range table.Paths {
				fmt.Printf("Received path: %+v\n", path)
				for i, attr := range path.Pattrs {
					fmt.Printf("Attribute %d: TypeUrl=%s\n", i, attr.TypeUrl)
					if attr.TypeUrl == "type.googleapis.com/gobgpapi.UnknownAttribute" {
						var unknownAttr api.UnknownAttribute
						if err := proto.Unmarshal(attr.Value, &unknownAttr); err != nil {
							log.Printf("Failed to unmarshal unknown attribute: %v", err)
							continue
						}
						fmt.Printf("Unknown Attribute Type: %d, Flags: %d, Value: %v\n", unknownAttr.Type, unknownAttr.Flags, unknownAttr.Value)
						if unknownAttr.Type == 99 {
							var config DeviceConfig
							if err := json.Unmarshal(unknownAttr.Value, &config); err != nil {
								log.Printf("Failed to unmarshal config: %v", err)
								continue
							}
							fmt.Printf("Received configuration:\n")
							fmt.Printf("Device ID: %s\n", config.DeviceID)
							fmt.Printf("Interface: %s\n", config.Config.Interface)
							fmt.Printf("IP Address: %s\n", config.Config.IPAddress)
							fmt.Printf("Netmask: %s\n", config.Config.Netmask)
						}
					}
				}
			}
		}
	}); err != nil {
		log.Fatal(err)
	}

	select {}
}
