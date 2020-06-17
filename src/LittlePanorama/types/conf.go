package types

import (
	pb "LittlePanorama/build/gen"
)

type HealthServerConfig struct{
	Addr string
	Id string
	Subjects []string
	Peers []*pb.Peer
}

func EmptyConf() *HealthServerConfig{
	return &HealthServerConfig{}
}

func SingleServerConf(addr,id string) *HealthServerConfig{
	peers := make([]*pb.Peer,1)
	peers[0] = &pb.Peer{Addr:addr, Id: id}
	return &HealthServerConfig{Addr:addr, Id:id, Peers:peers}
}