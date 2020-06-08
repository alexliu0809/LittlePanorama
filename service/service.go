package service
import (
	"fmt"
	"errors"
	"google.golang.org/grpc"
	"net"
	"golang.org/x/net/context"
	"LittlePanorama/store"
	"LittlePanorama/decision"
	"LittlePanorama/exchange"
	"github.com/golang/protobuf/ptypes"
	"LittlePanorama/types"

	rf "google.golang.org/grpc/reflection"
	pb "LittlePanorama/build/gen"
	
)

type PanoramaServer struct{
	me *pb.Peer // the addr and this server listens to and its id
	peers []*pb.Peer // the addr of its peers	
	
	// a unit responsible for commnuication
	port_listener net.Listener
	grpc_server *grpc.Server

	// status of the server
	status int

	// in memory storage unit
	storage store.Storage

	// exchanger
	exchanger exchange.Exchanger
}

// Initialize an instance of the panorama server
func NewPanoramaServer(conf *types.Conf) *PanoramaServer{
	me := &pb.Peer{Id:conf.Id,Addr:conf.Addr}
	peers := conf.Peers
	fmt.Printf("Configs:\nMe:%s\nPeers:%s\nSubjects:%s\n\n",me,peers,conf.Subjects)
	return &PanoramaServer{me:me, peers:peers, status:0, storage:store.NewStorage(), exchanger:exchange.NewExchanger(me,peers)}
}

// Start
func (self *PanoramaServer) Start() error{
	//fmt.Println("Starting Server")
	if self.status == 1{
		return errors.New("Server Already Running")
	}
	
	/******* Getting Ready to Server *******/
	// https://grpc.io/docs/tutorials/basic/go/
	// listen to a port
	var e error
	self.port_listener, e = net.Listen("tcp", self.me.Addr)
	if e != nil{
		return e
	}
	// start a new grpc server
	self.grpc_server = grpc.NewServer()
	// register with pb module
	pb.RegisterHealthServiceServer(self.grpc_server,self)
	// register with reflection
	rf.Register(self.grpc_server)
	// mark as online
	self.status = 1
	// get ready for serving, returns an error if necessary
	fmt.Println("Start Serving at Address:",self.me.Addr,"\n")
	return self.grpc_server.Serve(self.port_listener)
}

func (self *PanoramaServer) Stop() error{
	if self.status == 0{
		return errors.New("Server Already Stopped")
	}
	self.grpc_server.Stop()
	// reset everything
	self.grpc_server = nil
	self.port_listener = nil
	self.status = 0
	return nil
}

// Register a local observer to the health service.
// Must be called before SubmitReport.
func (self *PanoramaServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error){
	new_observer := store.Observer{Module:in.Module, Name:in.Observer}
	// id assigned to this observer. if multiple duplicate registrations, 
	// we would only assign the same handle
	observer_handle, err := self.storage.Register(new_observer)
	fmt.Println("Register",in.Module,in.Observer,"\n")
	return &pb.RegisterReply{Handle: observer_handle}, err
}

// Submit a report to the view storage
func (self *PanoramaServer) SubmitReport(ctx context.Context, in *pb.SubmitReportRequest) (*pb.SubmitReportReply, error){
	// ***** Note: todo: invalid health metric ****
	if self.storage.ValidHandle(in.Handle) == false{
		var result = pb.SubmitReportReply_FAILED
		return &pb.SubmitReportReply{Result: result}, errors.New("Invalid Handle")
	}
	var result = pb.SubmitReportReply_ACCEPTED
	fmt.Printf("Report Received on %s:\n%s\n\n",in.Report.Subject,in.Report.Observation)
	var is_observing bool = self.storage.IsObserving(in.Report.Subject)
	self.storage.SubmitReport(in.Report)
	
	
	// I am gonna observe this subject once I have one report
	if is_observing == false{
		go self.Observe(ctx, &pb.ObserveRequest{Subject:in.Report.Subject})
	}
	// Propagate this report to others who are observing it.
	go self.exchanger.PropageNewReport(in.Report)

	return &pb.SubmitReportReply{Result: result}, nil
}

// Query the latest raw health report of an entity
func (self *PanoramaServer) GetLatestReport(ctx context.Context, in *pb.GetReportRequest) (*pb.Report, error){
	r := self.storage.GetLatestReport(in.Subject)
	if r == nil{
		return nil, errors.New("No Report Available For "+in.Subject)
	}
	return r, nil
}
// Get all reports on one subject
func (self *PanoramaServer) GetPanorama(ctx context.Context, in *pb.GetPanoramaRequest) (*pb.Panorama, error){
	p := self.storage.GetPanorama(in.Subject)
	if p == nil{
		return nil, errors.New("No Panorama Available For "+in.Subject)
	}
	return p, nil
}	

// Get all reports from an observer on one subject
func (self *PanoramaServer) GetView(ctx context.Context, in *pb.GetViewRequest) (*pb.View, error){
	v := self.storage.GetView(in.Subject, in.Observer)
	if v == nil{
		return nil, errors.New("No View Available For Subject "+in.Subject+" From Observer "+in.Observer)
	}
	return v, nil
}

// Query a summarized health report from different observers about an entity
func (self *PanoramaServer) GetInference(ctx context.Context, in *pb.GetInferenceRequest) (*pb.Inference, error){
	p := self.storage.GetPanorama(in.Subject)
	if p == nil{
		return nil, errors.New("No Inference Available For "+in.Subject)
	}
	inference := decision.Inference(p)
	fmt.Println(inference.Observation.Metrics)
	return inference, nil
}

// Get all the peers of this DH server
func (self *PanoramaServer) GetPeers(ctx context.Context, in *pb.Empty) (*pb.GetPeerReply, error){
	fmt.Println("Get Peers")
	r := make([]*pb.Peer, len(self.peers)-1)
	for _, p := range self.peers{
		if p.Addr != self.me.Addr && p.Id != self.me.Id{
			r = append(r,p)
		}
	}
	return &pb.GetPeerReply{Peers:r},nil
}

// Get the ID of this health server
func (self *PanoramaServer) GetId(ctx context.Context, in *pb.Empty) (*pb.Peer, error){
	return self.me, nil
}

// Receive a ping from client
func (self *PanoramaServer) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error){
	fmt.Println("Receiving a ping from ",in.Source.Addr)
	t := ptypes.TimestampNow()
	return &pb.PingReply{Time: t, Result: pb.PingReply_GOOD}, nil
}

// Query the list of all subjects that have been observed
func (self *PanoramaServer) GetObservedSubjects(ctx context.Context, in *pb.Empty) (*pb.GetObservedSubjectsReply, error){
	fmt.Println("Get Observed Subjects")
	subjects := self.storage.DumpSubjects()
	if subjects == nil{
		return nil, errors.New("No Subject Available For Dump")
	}
	return &pb.GetObservedSubjectsReply{Subjects:subjects}, nil
}

// Dump all the raw health reports about all observed entities
func (self *PanoramaServer) DumpPanorama(ctx context.Context, in *pb.Empty) (*pb.DumpPanoramaReply, error){
	panoramas := self.storage.DumpPanorama()
	if panoramas == nil{
		return nil, errors.New("No Panorama Available For Dump")
	}
	return &pb.DumpPanoramaReply{Panoramas:panoramas}, nil
}

// Dump all the inferred health reports about all observed entities
func (self *PanoramaServer) DumpInference(ctx context.Context, in *pb.Empty) (*pb.DumpInferenceReply, error){
	panoramas := self.storage.DumpPanorama()
	if panoramas == nil{
		return nil, errors.New("No Inference Available For Dump")
	}
	inferences := make(map[string]*pb.Inference)
	for subject, panorama := range panoramas{
		inference := decision.Inference(panorama)
		inferences[subject] = inference
	}
	return &pb.DumpInferenceReply{Inferences:inferences}, nil
}

// Client requires you to observe a subject
// Now you should send this info to other peers
func (self *PanoramaServer) Observe(ctx context.Context, in *pb.ObserveRequest) (*pb.ObserveReply, error){
	reply, err := self.exchanger.Observe(in.Subject)
	return reply, err
}

// Stop observing a particular subject, all the reports
// concerning this subject will be ignored
func (self *PanoramaServer) StopObserving(ctx context.Context, in *pb.ObserveRequest) (*pb.ObserveReply, error){
	reply, err := self.exchanger.StopObserving(in.Subject)
	return reply, err
}

// Receive a report from a peer after you subscribe
func (self *PanoramaServer) LearnReport(ctx context.Context, in *pb.LearnReportRequest) (*pb.LearnReportReply, error){
	// request_type := LearnReportRequest.Kind
	// source := LearnReportRequest.Peer
	// report = LearnReportRequest.Report
	// handles different request types differently
	switch in.Kind{

    case pb.LearnReportRequest_SUBSCRIPTION:
        reply, err := self.exchanger.PeerObserve(in.Source, in.Report.Subject)
		return reply, err
	case pb.LearnReportRequest_UNSUBSCRIPTION:
        reply, err := self.exchanger.PeerUnobserve(in.Source, in.Report.Subject)
		return reply, err
	case pb.LearnReportRequest_NORMAL:
		fmt.Printf("New Report Received from %s(%s) on %s\n\n",in.Source.Id,in.Source.Addr,in.Report.Subject)
		reply, err := self.addReport(in.Report)
		return reply, err
    }

	return nil, nil
}

func (self *PanoramaServer) addReport(report *pb.Report) (*pb.LearnReportReply, error){
	self.storage.SubmitReport(report)
	return &pb.LearnReportReply{Result:pb.LearnReportReply_ACCEPTED}, nil
}