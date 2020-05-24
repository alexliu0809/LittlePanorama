package main

import (
	"LittlePanorama/service"
	pb "LittlePanorama/build/gen"
)

func main(){
	sever := service.NewPanoramaServer(&pb.Peer{Id:"pano0",Addr:"127.0.0.1:6688"},nil)
	sever.Start()
}