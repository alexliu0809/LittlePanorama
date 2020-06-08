package main

import (
	"flag"
	"strings"
	"fmt"
	"LittlePanorama/service"
	"LittlePanorama/types"
	ut "LittlePanorama/util"
	//pb "LittlePanorama/build/gen"
)

func main(){
	/*** Flags ***/
	var config     = flag.String("config", "", "read config file that contains peer addr")
	var addr       = flag.String("addr", "localhost", "this server listens to")

	var conf *types.Conf
	// parse the flag
	flag.Parse()
	// parse rc
	conf, err := ut.ParseRC(*config)
	if err != nil{
		fmt.Println(err)
		fmt.Println("--------- Usage: ---------")
		flag.PrintDefaults()
		return
	} else if conf == nil{
		// no conf, read addr
		me_addr := *addr
		if !strings.ContainsAny(me_addr, ":") {
			fmt.Println("Error: Wrong Format of Addr Specified")
			fmt.Println("--------- Usage: ---------")
			flag.PrintDefaults()
			return
		}
		if len(flag.Args()) == 0 {
			fmt.Println("Error: No Server Id Specified")
			fmt.Println("--------- Usage: ---------")
			flag.PrintDefaults()
			return
		} else if len(flag.Args()) > 1{
			fmt.Println("Error: Too Many Args")
			fmt.Println("--------- Usage: ---------")
			flag.PrintDefaults()
			return
		}
		me_id := flag.Args()[0]
		conf = types.SingleServerConf(me_addr, me_id)
	}

	sever := service.NewPanoramaServer(conf)
	sever.Start()
}
