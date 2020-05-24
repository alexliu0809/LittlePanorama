package decision

import (
//	"fmt"
	pb "LittlePanorama/build/gen"
	ts "github.com/golang/protobuf/ptypes/timestamp"
)

// infer the status of a subject given its panorama
func Inference(in *pb.Panorama) *pb.Inference{
	if in == nil{ // no panorama is available
		return nil
	}

	var max_ts *ts.Timestamp = nil
	observers_map := make(map[string]bool)
	score_of_status_of_metric := make(map[string]map[pb.Status]float32)
	count_of_status_of_metric := make(map[string]map[pb.Status]int)
	// majority vote from all observers for each metric
	// iterate all the views
	for _,view := range in.Views{
		// first aggragate the results for one view
		ob_from_one_view := infer_one_view(view)
		// track all observers
		observers_map[view.Observer] = true
		// track ts
		if max_ts == nil || ts_cmp(ob_from_one_view.Ts,max_ts) > 0{
			max_ts = ob_from_one_view.Ts
		}
		// for each metric in that view
		for _, metric := range ob_from_one_view.Metrics{
			metric_name := metric.Name
			metric_status := metric.Value.Status
			metric_score := metric.Value.Score
			_, found := count_of_status_of_metric[metric_name]
			if found == false{
				status_score_map := make(map[pb.Status]float32)
				status_count_map := make(map[pb.Status]int)
				status_score_map[metric_status] = metric_score
				status_count_map[metric_status] = 1
				// initialize for a new metric
				score_of_status_of_metric[metric_name] = status_score_map
				count_of_status_of_metric[metric_name] = status_count_map
			} else {
				// metric already exists.
				// but status might not exist in the inner map
				status_score_map := score_of_status_of_metric[metric_name]
				status_count_map := count_of_status_of_metric[metric_name]
				_, found = status_score_map[metric_status]
				// this status doesn't exist in the map 
				if found == false{
					status_score_map[metric_status] = metric_score
					status_count_map[metric_status] = 1
				} else {
					status_score_map[metric_status] += metric_score
					status_count_map[metric_status] += 1
				}
			}
		}
	}
	// done with all stats.
	// summarize with majority vote
	// convert all the stats to one observation
	metrics := make(map[string]*pb.Metric,0)
	for metric_name := range count_of_status_of_metric{
		status_score_map := score_of_status_of_metric[metric_name]
		status_count_map := count_of_status_of_metric[metric_name]
		max_cnt := -1000
		var max_score float32
		var max_status pb.Status
		for status_name := range status_count_map{
			status_count := status_count_map[status_name]
			status_score := status_score_map[status_name]
			if status_count > max_cnt {
				max_cnt = status_count
				max_status = status_name
				max_score = status_score
			}
		}
		avg_max_score := max_score / float32(max_cnt)
		metric_value := &pb.Value{Status:max_status, Score:avg_max_score}
		metrics[metric_name] = &pb.Metric{Name:metric_name,Value:metric_value}
	}
	// an observation that summarizes everything
	ttl_ob := &pb.Observation{Ts:max_ts, Metrics:metrics}

	// covert observers to []string
	observers := make([]string,len(observers_map))
	for k := range observers_map{
		observers = append(observers,k)
	}

	return &pb.Inference{Subject:in.Subject,Observers:observers,Observation:ttl_ob}
}

func infer_one_view(in *pb.View) *pb.Observation{
	if in == nil{
		return nil // no view is available
	}

	// stop when we hit a status that is different than the newest status
	// keep track of all the stats
	score_of_metric := make(map[string]float32)
	count_of_metric := make(map[string]int)
	newest_status_of_metric := make(map[string]pb.Status)
	end_of_metric := make(map[string]bool)

	// iterate observations from the most up-to-date to least up-to-date
	for i := len(in.Observations) - 1; i >= 0; i --{
		ob := in.Observations[i]
		for _, metric := range ob.Metrics{
			metric_name := metric.Name
			metric_status := metric.Value.Status
			metric_score := metric.Value.Score
			// get the newest status
			newest, found := newest_status_of_metric[metric_name]
			// not found, you are the newest
			if found == false{
				newest_status_of_metric[metric_name] = metric_status
				count_of_metric[metric_name] = 1
				score_of_metric[metric_name] = metric_score
				end_of_metric[metric_name] = false
			} else {
				// we are not recording this metric anymore
				if end_of_metric[metric_name] == true{
					continue
				} else {
					// if this metric doesn't agree with the newest one, stop recording
					if two_status_the_same(newest, metric_status) == false{
						end_of_metric[metric_name] = true
						continue
					} else {
						// keep track of all scores
						count_of_metric[metric_name] += 1
						score_of_metric[metric_name] += metric_score
					}
				}
			}
		}
	}
	// convert all the stats to one observation
	metrics := make(map[string]*pb.Metric,0)
	for metric_name := range newest_status_of_metric{
		avg_score := score_of_metric[metric_name] / float32(count_of_metric[metric_name])
		metric_value := &pb.Value{Status:newest_status_of_metric[metric_name], Score:avg_score}
		metrics[metric_name] = &pb.Metric{Name:metric_name,Value:metric_value}
	}

	return &pb.Observation{Ts:in.Observations[len(in.Observations)-1].Ts, Metrics:metrics}
}

// compare two status are the same. 
func two_status_the_same(newest pb.Status, this pb.Status) bool{
	// if the newest is HEALTHY, it is compatible with pending
	if newest == pb.Status_HEALTHY{
		if this == pb.Status_HEALTHY || this == pb.Status_PENDING{
			return true
		} else {
			return false
		}
	} else {
		return newest == this
	}
}

func ts_cmp(one *ts.Timestamp, another *ts.Timestamp) int32 {
	if one.Seconds == another.Seconds {
	 	return one.Nanos - another.Nanos
	} else {
		if one.Seconds > one.Seconds{
			return 1
		} else {
			return -1
		}
	}
}