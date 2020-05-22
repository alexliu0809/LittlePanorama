package service
import (
	"fmt"
	"errors"
	"google.golang.org/grpc"
	"net"
	"golang.org/x/net/context"
	rf "google.golang.org/grpc/reflection"
	pb "LittlePanorama/build/gen"
)

type PanoramaServer struct{
	addr string // the addr this server listens to 
	peers []string // the addr of its peers	
	
	// a unit responsible for commnuication
	port_listener net.Listener
	grpc_server *grpc.Server

	// status of the server
	status int
}




// Initialize an instance of the panorama server
func NewPanoramaServer(addr string, peers []string) *PanoramaServer{
	return &PanoramaServer{addr:addr, peers:peers}
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
	pb.RegisterPanoramaServiceServer(self.grpc_server,self)
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

func (self *PanoramaServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error){
	return &pb.RegisterReply{Handle: 0}, nil
}


