/**** Created by Alex Liu to compare LittlePanorama and Panorama ****/
package service

import (
	"testing"
	"golang.org/x/net/context"
	"time"
	"sync"
	"math/rand"
	"panorama/types"
	"google.golang.org/grpc"
	"fmt"
	"os"
	pt "github.com/golang/protobuf/ptypes"
	pb "panorama/build/gen"
)

var client pb.HealthServiceClient
var handle uint64
var peers []pb.HealthServiceClient
var peer_addrs []string
var r = rand.New(rand.NewSource(int64(0)))
var test_subject string = "test_subject"
var test_observer string = "test_observer"
var rand_statuses = [2]pb.Status{pb.Status_HEALTHY,pb.Status_PENDING}
var all_names = [3]string{"cpu","memory","network"}

/***** Helper Functions *****/
func randReport() *pb.Report{
	return createReport(test_observer,test_subject,randMetrics())
}

func randMetrics() map[string]*pb.Metric{
	names := make([]string,3)
	statuses := make([]pb.Status,3)
	scores := make([]float32,3)
	// initialize 
	for i := 0; i<3; i++{
		names[i] = all_names[i]
		statuses[i] = rand_statuses[rand.Int()%len(rand_statuses)]
		scores[i] = float32(int(rand.Float32() * 10))
	}
	return createMetrics(names,statuses,scores)
}

func createReport(observer string, subject string, metrics map[string]*pb.Metric) *pb.Report{
	return &pb.Report{Observer:observer,Subject:subject,Observation:createObservation(metrics)}
}

func createMetrics(names []string, statuses []pb.Status, scores []float32) map[string]*pb.Metric{
	metrics := make(map[string]*pb.Metric)
	for i := 0; i < len(names); i++{
		this_name := names[i]
		this_status := statuses[i]
		this_score := scores[i]
		metrics[this_name] = &pb.Metric{Name:this_name,Value: &pb.Value{Status:this_status,Score:this_score}}
	}
	return metrics
}

func createObservation(metrics map[string]*pb.Metric) *pb.Observation{
	ts, _ := pt.TimestampProto(time.Now())
	return &pb.Observation{Ts: ts, Metrics: metrics}
}


/***** Benchmark the speed of submitting report *****/
func BenchmarkSubmitReport(b *testing.B){
	r := randReport()
	for i := 0; i < b.N; i++{
		client.SubmitReport(context.Background(), &pb.SubmitReportRequest{Handle: handle, Report: r})
	}
}

/***** A base function that allows you to vary the number of reports *****/
func benchmarkGetInference(num_reports int, b *testing.B){
	// submit n reports
	for i := 0; i < num_reports; i++{
		r := randReport()
		client.SubmitReport(context.Background(), &pb.SubmitReportRequest{Handle: handle, Report: r})
	}
	// get inference
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.GetInference(context.Background(), &pb.GetInferenceRequest{Subject:test_subject})
	}
}

func BenchmarkGetInference100(b *testing.B)  { benchmarkGetInference(100, b) }
//func BenchmarkGetInference200(b *testing.B)  { benchmarkGetInference(200, b) }
//func BenchmarkGetInference400(b *testing.B)  { benchmarkGetInference(400, b) }
// func BenchmarkGetInference800(b *testing.B) { benchmarkGetInference(800, b) }
// func BenchmarkGetInference1600(b *testing.B) { benchmarkGetInference(1600, b) }
// func BenchmarkGetInference3200(b *testing.B) { benchmarkGetInference(3200, b) }

/***** A base function that allows you to vary the number of clients *****/
func benchmarkPropagate(num_peers int ,b *testing.B){
	r := randReport()
	b.ResetTimer()
	for i := 0; i < b.N; i++{
		request := &pb.LearnReportRequest{Source: &pb.Peer{Addr:"127.0.0.1",Id:"client"}, Report: r}
		var wg sync.WaitGroup
		for j := 0; j < num_peers; j++ {
			peer := peers[j]
			wg.Add(1)
			go func() {
				_, err := peer.LearnReport(context.Background(), request)
				if err != nil {
					fmt.Errorf("failed to propagate report about %s to %s\n", r.Subject, peer)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkPropagate2(b *testing.B)  { benchmarkPropagate(2, b) }
func BenchmarkPropagate4(b *testing.B)  { benchmarkPropagate(4, b) }
func BenchmarkPropagate8(b *testing.B)  { benchmarkPropagate(8, b) }
func BenchmarkPropagate16(b *testing.B)  { benchmarkPropagate(16, b) }
func BenchmarkPropagate32(b *testing.B)  { benchmarkPropagate(32, b) }

func TestMain(m *testing.M) {
	// set up one server
	config := &types.HealthServerConfig{ Addr:"127.0.0.1:6688", Id:"pano0"}
	server := NewHealthGServer(config)
	go server.Start(make(chan error))
	time.Sleep(1 * time.Second)

	// create a single client for testing
	conn, _ := grpc.Dial("127.0.0.1:6688", grpc.WithInsecure())
	defer conn.Close()
	client = pb.NewHealthServiceClient(conn)
	reply, _ := client.Register(context.Background(), &pb.RegisterRequest{Module: "Test", Observer: test_observer})
	handle = reply.Handle
	
	// create multiple clients
	max_num_peers := 32
	peer_addrs = make([]string,max_num_peers)
	peers = make([]pb.HealthServiceClient,max_num_peers)
	for i := 0; i < max_num_peers; i ++{
		peer_addrs[i] = fmt.Sprintf("127.0.0.1:%d",6690+i)
	}

	for i := 0; i < max_num_peers; i ++{
		conn, err := grpc.Dial(peer_addrs[i], grpc.WithInsecure())
		if err != nil {
			fmt.Errorf("Failed to connect to peerat %s\n", peer_addrs[i])
		}
		peers[i]=pb.NewHealthServiceClient(conn)
	}
	time.Sleep(1 * time.Second)
	
	// run benchmarks
	exitVal := m.Run()
	os.Exit(exitVal)
}
