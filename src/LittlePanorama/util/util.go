package util
import (
	"LittlePanorama/types"
	"os"
	"encoding/json"
	"io/ioutil"
	"errors"
	"strings"
	//"fmt"
	pb "LittlePanorama/build/gen"
)

func ParseRC(rc string) (*types.HealthServerConfig, error) {
	if len(rc) > 0{

		f, err := os.Open(rc)
		defer f.Close()
		if err != nil{
			return nil, err
		}
		byteValue, err := ioutil.ReadAll(f)
		if err != nil{
			return nil, err
		}
		var result map[string]interface{}
		json.Unmarshal([]byte(byteValue), &result)
		// contruct conf
		conf := types.EmptyConf()
		if addr, found := result["Addr"]; found == false{
			return nil, errors.New("Not Addr For This Server Specified in Config")
		} else {
			conf.Addr, _ = addr.(string)
		}

		if id, found := result["Id"]; found == false{
			return nil, errors.New("Not Id For This Server Specified in Config")
		} else {
			conf.Id, _ = id.(string)
		}

		if subjects, found := result["Subjects"]; found == true{
			conf.Subjects = subjects.([]string)
		}

		if raw_peers, found := result["Peers"]; found == false{
			return nil, errors.New("Not Peers For This Server Specified in Config")
		} else {
			parsed_peers, _ := raw_peers.(map[string]interface{})
			//fmt.Println(parsed_peers, err)
			pb_peers := make([]*pb.Peer,0)
			for peer_id, peer_addr := range parsed_peers{
				if !strings.ContainsAny(peer_addr.(string), ":"){
					return nil, errors.New("Peer Addr in Wrong Format in Config")
				} else {
					pb_peers = append(pb_peers, &pb.Peer{Id:peer_id,Addr:peer_addr.(string)})
				}
			}
			
			conf.Peers = pb_peers
			return conf, nil
		}
	} else {
		return nil, nil
	}
	
}