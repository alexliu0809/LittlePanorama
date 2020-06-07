package exchange

import (
	pb "LittlePanorama/build/gen"
	"golang.org/x/net/context"
	//"LittlePanorama/store"
	"errors"
)

type Exchanger struct{
	me *pb.Peer // the addr and this server listens to and its id
	peers []*pb.Peer // the addr of its peers

	// a peer could subscribe to a list of subjects
	peers_subscribe_subjects map[*pb.Peer][]string

	// a subject could be subscribed by multipled objects
	subjects_subscribed_by_peers map[string][]*pb.Peer

	// saved connections for each peer
	connections map[*pb.Peer]pb.HealthServiceClient
}

func NewExchanger(me *pb.Peer, peers []*pb.Peer) Exchanger{
	return Exchanger{me, peers, make(map[*pb.Peer][]string), make(map[string][]*pb.Peer), make(map[*pb.Peer]pb.HealthServiceClient)}
}

// send learn report request to all peers about subscription
// use the learn report function to achieve this
func (self *Exchanger) Subscribe(ctx context.Context, subject string) error{
	report := &pb.Report{Subject:subject}
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_SUBSCRIPTION,Source:self.me,Report:report}
	for _, peer := range self.peers{
		conn, err := self.connect(peer)
		if err != nil{
			return err
		}
		reply, err := conn.LearnReport(ctx, request)
		if err != nil{
			return err
		} else if reply.Result == pb.LearnReportReply_FAILED {
			return errors.New("Learn Report Failed")
		}
	}
	return nil
}

func (self *Exchanger) Unsubscribe(ctx context.Context, subject string) error{
	return nil
}

func (self *Exchanger) connect(peer *pb.Peer) (pb.HealthServiceClient, error){
	return nil, nil
}

func (self *Exchanger) NotifyPeers(in *pb.LearnReportRequest) error{
	return nil
}


