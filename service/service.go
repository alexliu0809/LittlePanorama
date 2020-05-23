package service
import (
	"fmt"
	"errors"
	"google.golang.org/grpc"
	"net"
	"golang.org/x/net/context"
	rf "google.golang.org/grpc/reflection"
	pb "LittlePanorama/build/gen"
	"LittlePanorama/store"
)

type PanoramaServer struct{
	addr string // the addr this server listens to 
	peers []string // the addr of its peers	
	
	// a unit responsible for commnuication
	port_listener net.Listener
	grpc_server *grpc.Server

	// status of the server
	status int

	// in memory storage unit
	storage store.Storage
}


// Initialize an instance of the panorama server
func NewPanoramaServer(addr string, peers []string) *PanoramaServer{
	return &PanoramaServer{addr:addr, peers:peers, status:0, storage:store.NewStorage()}
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
	self.port_listener, e = net.Listen("tcp", self.addr)
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
	fmt.Println("Start Serving at Address:",self.addr)
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
	fmt.Println("Register",in.Module,in.Observer)
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
	fmt.Println(in.Report.Subject)
	fmt.Println(in.Report.Observer)
	self.storage.SubmitReport(in.Report)
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
	return nil, nil
}
// Get all the peers of this DH server
func (self *PanoramaServer) GetPeers(ctx context.Context, in *pb.Empty) (*pb.GetPeerReply, error){
	return nil, nil
}
// Get the ID of this health server
func (self *PanoramaServer) GetServerId(ctx context.Context, in *pb.Empty) (*pb.Peer, error){
	return nil, nil
}


