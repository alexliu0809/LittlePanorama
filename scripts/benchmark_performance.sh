#!/bin/bash

# scripts to benchmark the performance of both LittlePanorama and Panorama

cd /home/ubuntu/go/src/panorama/service
echo "panorama BenchmarkSubmitReport"
go test -bench=BenchmarkSubmitReport | grep "ns/op"
cd /home/ubuntu/go/src/LittlePanorama/service
echo "LittlePanorama BenchmarkSubmitReport"
go test -bench=BenchmarkSubmitReport | grep "ns/op"



cd /home/ubuntu/go/src/panorama/service
echo "panorama BenchmarkInference"
go test -bench=BenchmarkGetInference | grep "ns/op"
cd /home/ubuntu/go/src/LittlePanorama/service
echo "LittlePanorama BenchmarkInference"
go test -bench=BenchmarkGetInference | grep "ns/op"


cd /home/ubuntu/go/src/panorama/service
echo "panorama BenchmarkPropagate"
go test -bench=BenchmarkPropagate | grep "ns/op"
cd /home/ubuntu/go/src/LittlePanorama/service
echo "LittlePanorama BenchmarkPropagate"
go test -bench=BenchmarkPropagate | grep "ns/op"

