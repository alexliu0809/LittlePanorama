package exchange

import (
	pb "LittlePanorama/build/gen"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

// send learn report request to all peers about observing this subject
// use the learn report function to achieve this
func (self *Exchanger) Observe(ctx context.Context, subject string) (*pb.ObserveReply, error){
	report := &pb.Report{Subject:subject}
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_SUBSCRIPTION,Source:self.me,Report:report}
	err := self.NotifyPeers(ctx, request)
	if err == nil{
		return &pb.ObserveReply{Success:true},nil
	} else {
		return &pb.ObserveReply{Success:false},err
	}
}

func (self *Exchanger) StopObserving(ctx context.Context, subject string) (*pb.ObserveReply, error){
	report := &pb.Report{Subject:subject}
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_UNSUBSCRIPTION,Source:self.me,Report:report}
	err := self.NotifyPeers(ctx, request)
	if err == nil{
		return &pb.ObserveReply{Success:true},nil
	} else {
		return &pb.ObserveReply{Success:false},err
	}
}

func (self *Exchanger) PeerObserve(peer *pb.Peer, subject string) (*pb.LearnReportReply, error){
	if peer.Addr == self.me.Addr && peer.Id == self.me.Id{
		return &pb.LearnReportReply{Result:pb.LearnReportReply_FAILED},errors.New("Peer is Yourself")
	}

	for _, subscribed := range self.peers_subscribe_subjects[peer]{
		if subscribed == subject{
			return &pb.LearnReportReply{Result:pb.LearnReportReply_IGNORED}, nil
		}
	}
	self.peers_subscribe_subjects[peer] = append(self.peers_subscribe_subjects[peer], subject)
	return &pb.LearnReportReply{Result:pb.LearnReportReply_ACCEPTED}, nil
}

func (self *Exchanger) PeerUnobserve(peer *pb.Peer, subject string) (*pb.LearnReportReply, error){
	if peer.Addr == self.me.Addr && peer.Id == self.me.Id{
		return nil, errors.New("Peer is Yourself")
	}

	for index, subscribed := range self.peers_subscribe_subjects[peer]{
		if subscribed == subject{
			self.peers_subscribe_subjects[peer] = append(self.peers_subscribe_subjects[peer][:index], self.peers_subscribe_subjects[peer][index+1:]...)
			return &pb.LearnReportReply{Result:pb.LearnReportReply_ACCEPTED}, nil
		}
	}
	return &pb.LearnReportReply{Result:pb.LearnReportReply_IGNORED}, nil
}

func (self *Exchanger) connect(peer *pb.Peer) (pb.HealthServiceClient, error){
	if conn, found := self.connections[peer]; found == false{
		new_conn, err := grpc.Dial(peer.Addr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		self.connections[peer] = pb.NewHealthServiceClient(new_conn)
		return self.connections[peer], nil
	} else {
		return conn, nil
	}
}

func (self *Exchanger) NotifyPeers(ctx context.Context, request *pb.LearnReportRequest) error{
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

func (self *Exchanger) PropageNewReport(ctx context.Context, report *pb.Report){
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_NORMAL,Source:self.me,Report:report}
	self.NotifyPeers(ctx, request)
}