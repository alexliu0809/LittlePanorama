package exchange

import (
	pb "LittlePanorama/build/gen"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"LittlePanorama/store"
	"fmt"
	"errors"
)

type Exchanger struct{
	me *pb.Peer // the addr and this server listens to and its id
	peers []*pb.Peer // the addr of its peers

	// a peer could subscribe to a list of subjects
	// peers_subscribe_subjects map[*pb.Peer][]string

	// a subject could be subscribed by multipled objects
	subjects_subscribed_by_peers map[string][]*pb.Peer

	// saved connections for each peer
	connections map[*pb.Peer]pb.HealthServiceClient
}

func NewExchanger(me *pb.Peer, peers []*pb.Peer) Exchanger{
	return Exchanger{me, peers, make(map[string][]*pb.Peer), make(map[*pb.Peer]pb.HealthServiceClient)}
}

// send learn report request to all peers about observing this subject
// use the learn report function to achieve this
func (self *Exchanger) Observe(subject string) (*pb.ObserveReply, error){
	report := &pb.Report{Subject:subject}
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_SUBSCRIPTION,Source:self.me,Report:report}
	fmt.Println("Start Observing",subject,"\n")
	err := self.NotifyPeers(request)
	if err == nil{
		return &pb.ObserveReply{Success:true},nil
	} else {
		return &pb.ObserveReply{Success:false},err
	}
}

func (self *Exchanger) StopObserving(subject string) (*pb.ObserveReply, error){
	report := &pb.Report{Subject:subject}
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_UNSUBSCRIPTION,Source:self.me,Report:report}
	fmt.Println("Stop Observing ",subject,"\n")
	err := self.NotifyPeers(request)
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
	fmt.Printf("Peer %s at %s starts observing %s\n\n", peer.Id, peer.Addr, subject)
	if peers, found := self.subjects_subscribed_by_peers[subject]; found == false{
		peers = make([]*pb.Peer,1)
		peers[0] = peer
		self.subjects_subscribed_by_peers[subject] = peers
	} else {
		for _, p := range peers{
			if p.Id == peer.Id && p.Addr == peer.Addr{
				return  &pb.LearnReportReply{Result:pb.LearnReportReply_IGNORED}, nil
			}
		}
		self.subjects_subscribed_by_peers[subject] = append(self.subjects_subscribed_by_peers[subject], peer)
	}
	return &pb.LearnReportReply{Result:pb.LearnReportReply_ACCEPTED}, nil
}

func (self *Exchanger) PeerUnobserve(peer *pb.Peer, subject string) (*pb.LearnReportReply, error){
	if peer.Addr == self.me.Addr && peer.Id == self.me.Id{
		return &pb.LearnReportReply{Result:pb.LearnReportReply_FAILED},errors.New("Peer is Yourself")
	}
	fmt.Printf("Peer %s at %s stops observing %s\n\n", peer.Id, peer.Addr, subject)
	if peers, found := self.subjects_subscribed_by_peers[subject]; found == false{
		return  &pb.LearnReportReply{Result:pb.LearnReportReply_IGNORED}, nil
	} else {
		for index, p := range peers{
			if p.Id == peer.Id && p.Addr == peer.Addr{
				// remove
				self.subjects_subscribed_by_peers[subject] = append(self.subjects_subscribed_by_peers[subject][:index],self.subjects_subscribed_by_peers[subject][:index+1]...)
				return  &pb.LearnReportReply{Result:pb.LearnReportReply_ACCEPTED}, nil
			}
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

func (self *Exchanger) NotifyPeers(request *pb.LearnReportRequest) error{
	for _, peer := range self.peers{
		// don't propage to yourself
		if peer.Id == self.me.Id && peer.Addr == self.me.Addr{
			continue
		}
		//fmt.Println("Connect")
		conn, err := self.connect(peer)
		if err != nil{
			return err
		}
		//fmt.Println("Learn Report")
		reply, err := conn.LearnReport(context.Background(), request)
		//fmt.Println(reply,err)
		if err != nil{
			return err
		} else if reply.Result == pb.LearnReportReply_FAILED {
			return errors.New("Learn Report Failed")
		}
	}
	return nil
}

func (self *Exchanger) PropageNewReport(report *pb.Report) error{
	request := &pb.LearnReportRequest{Kind:pb.LearnReportRequest_NORMAL,Source:self.me,Report:report}
	cnt := 0
	for _, peer := range self.subjects_subscribed_by_peers[report.Subject]{
		// don't propage to yourself
		if peer.Id == self.me.Id && peer.Addr == self.me.Addr{
			continue
		}
		conn, err := self.connect(peer)
		if err != nil{
			return err
		}
		fmt.Printf("Propagating Report to Peer %s at %s\n\n", peer.Id, peer.Addr)
		reply, err := conn.LearnReport(context.Background(), request)
		if err != nil{
			return err
		} else if reply.Result == pb.LearnReportReply_FAILED {
			return errors.New("Learn Report Failed")
		}
		cnt += 1
	}
	fmt.Printf("Report on %s propagated to %d peers\n\n", report.Subject, cnt)
	return nil
}