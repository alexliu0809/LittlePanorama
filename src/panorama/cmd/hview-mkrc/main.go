package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	dt "panorama/types"
	du "panorama/util"
)

var (
	tag       = "hview-mkrc"
	id        = flag.String("id", "", "name id of this server, must be one of the peers")
	nserver   = flag.Int("nserver", 3, "number of hview server")
	localhost = flag.Bool("localhost", true, "whether all servers are localhost")
	addressp  = flag.String("addressp", "", "pattern of server address, e.g., 10.10.2.%d")
	namep     = flag.String("namep", "DHS_%d", "pattern of server name, e.g., DHS_%d")
	namestart = flag.Int("namestart", 0, "start of server name id")
	sidstart  = flag.Int("sidstart", 0, "start of the server id in server pattern")
	portstart = flag.Int("port_start", 10000, "start of port range for a random port")
	portend   = flag.Int("port_end", 30000, "end of port range for a random port")
	fixport   = flag.Int("fix_port", -1, "fix port instead of random port number")
	subjects  = flag.String("subjects", "", "comma separated list of subjects to watch for health,\neffective only when FilterSubmission is true")
	filter    = flag.Bool("filter", false, "whether to filter health reports based on subjects")
	dbfile    = flag.String("dbfile", "deephealth.db", "database file to persist health information")
	output    = flag.String("output", "", "file path to output the generated RC")
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	flag.Parse()
	var p, s int
	var inc_p, inc_s bool
	if len(*id) == 0 {
		du.LogF(tag, "Must specify the name id of this server")
	}
	if *fixport > 0 {
		p = *fixport
		inc_p = false
	} else {
		if *portstart <= 0 || *portend <= 0 {
			du.LogF(tag, "Port range must be positive")
		}
		if *portstart > *portend {
			du.LogF(tag, "Port start must not exceed port end")
		}
		p = *portstart + int(r.Intn(*portend-*portstart))
		inc_p = true
	}
	if *sidstart < 0 {
		du.LogF(tag, "Server id must be positive")
	}
	s = *sidstart
	inc_s = len(*addressp) > 0

	rc := new(dt.HealthServerConfig)
	rc.Peers = make(map[string]string)
	for i := 0; i < *nserver; i++ {
		eid := fmt.Sprintf(*namep, i+*namestart)
		if len(*addressp) > 0 {
			rc.Peers[eid] = fmt.Sprintf(*addressp+":%d", s, p)
		} else {
			rc.Peers[eid] = fmt.Sprintf("localhost:%d", p)
		}
		if inc_p {
			p++
		}
		if inc_s {
			s++
		}
	}
	addr, ok := rc.Peers[*id]
	if !ok {
		du.LogF(tag, "%s is not one of the peers %v", *id, rc.Peers)
	}
	rc.Id = *id
	rc.Addr = addr
	rc.FilterSubmission = *filter
	if len(*subjects) > 0 {
		rc.Subjects = strings.Split(*subjects, ",")
	}
	rc.DBFile = *dbfile
	if len(*output) > 0 {
		fmt.Println("Saving to " + *output)
		err := dt.SaveConfig(*output, rc)
		if err != nil {
			du.LogF(tag, "Failed to save config: %v", err)
		}
		fmt.Println("Saved")
	} else {
		fmt.Println(dt.JString(rc))
	}
}
