package types

import (
	pb "LittlePanorama/build/gen"
)

type Conf struct{
	Addr string
	Id string
	Subjects []string
	Peers []*pb.Peer
}

func EmptyConf() *Conf{
	return &Conf{}
}

func SingleServerConf(addr,id string) *Conf{
	peers := make([]*pb.Peer,1)
	peers[0] = &pb.Peer{Addr:addr, Id: id}
	return &Conf{Addr:addr, Id:id, Peers:peers}
}