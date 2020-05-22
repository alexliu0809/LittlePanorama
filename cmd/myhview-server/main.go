package main

import (
	"LittlePanorama/service"
)

func main(){
	sever := service.NewPanoramaServer("127.0.0.1:6688",nil)
	sever.Start()	
}