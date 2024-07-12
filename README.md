# bgpfig

This is a simple and dumb experiment to play with the concept of sending configurations over BGP. This could simply
mean that we can send well known configuration items across via BGP.




```bash
➜  bgpfig git:(main) docker-compose up
WARN[0000] /Users/henryjc/GolandProjects/bgpfig/docker-compose.yml: `version` is obsolete
[+] Running 2/0
 ✔ Container bgpfig-bgp-server-1  Created                                                                                                                                                                                                0.0s
 ✔ Container bgpfig-bgp-client-1  Created                                                                                                                                                                                                0.0s
Attaching to bgp-client-1, bgp-server-1
bgp-server-1  | time="2024-07-12T00:45:56Z" level=info msg="Add a peer configuration" Key=172.20.0.3 Topic=Peer
bgp-server-1  | 2024/07/12 00:45:56 Sending path update with attributes: [[type.googleapis.com/apipb.OriginAttribute]:{} [type.googleapis.com/apipb.NextHopAttribute]:{next_hop:"172.20.0.2"} [type.googleapis.com/apipb.UnknownAttribute]:{flags:192  type:99  value:"{\"device_id\":\"router1\",\"config\":{\"interface\":\"GigabitEthernet0/1\",\"ip_address\":\"192.168.1.1\",\"netmask\":\"255.255.255.0\"}}"}]
bgp-server-1  | BGP speaker is running...
bgp-client-1  | time="2024-07-12T00:45:56Z" level=info msg="Add a peer configuration" Key=172.20.0.2 Topic=Peer
bgp-client-1  | time="2024-07-12T00:46:04Z" level=info msg="Peer Up" Key=172.20.0.2 State=BGP_FSM_OPENCONFIRM Topic=Peer
bgp-server-1  | time="2024-07-12T00:46:04Z" level=info msg="Peer Up" Key=172.20.0.3 State=BGP_FSM_OPENCONFIRM Topic=Peer
bgp-client-1  | Received path: nlri:{[type.googleapis.com/apipb.IPAddressPrefix]:{prefix_len:24 prefix:"192.168.1.0"}} pattrs:{[type.googleapis.com/apipb.OriginAttribute]:{}} pattrs:{[type.googleapis.com/apipb.AsPathAttribute]:{segments:{type:AS_SEQUENCE numbers:65003}}} pattrs:{[type.googleapis.com/apipb.NextHopAttribute]:{next_hop:"172.20.0.2"}} pattrs:{[type.googleapis.com/apipb.UnknownAttribute]:{flags:192 type:99 value:"{\"device_id\":\"router1\",\"config\":{\"interface\":\"GigabitEthernet0/1\",\"ip_address\":\"192.168.1.1\",\"netmask\":\"255.255.255.0\"}}"}} age:{seconds:1720745164} validation:{} family:{afi:AFI_IP safi:SAFI_UNICAST} source_asn:65003 source_id:"172.20.0.2" neighbor_ip:"172.20.0.2" local_identifier:1
bgp-client-1  | Attribute 0: TypeUrl=type.googleapis.com/apipb.OriginAttribute
bgp-client-1  | Attribute 1: TypeUrl=type.googleapis.com/apipb.AsPathAttribute
bgp-client-1  | Attribute 2: TypeUrl=type.googleapis.com/apipb.NextHopAttribute
bgp-client-1  | Attribute 3: TypeUrl=type.googleapis.com/apipb.UnknownAttribute
^CGracefully stopping... (press Ctrl+C again to force)
[+] Stopping 2/2
 ✔ Container bgpfig-bgp-client-1  Stopped                                                                                                                                                                                                0.1s
 ✔ Container bgpfig-bgp-server-1  Stopped

```