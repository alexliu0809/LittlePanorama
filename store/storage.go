package store

import (
	ts "github.com/golang/protobuf/ptypes/timestamp"
	pb "LittlePanorama/build/gen"
	//"LittlePanorama/decision"
	"sync"
	//"fmt"
)

type Observer struct{
	// which module this observer belongs to 
	Module string
	// name of the observer
	Name string
}

type Storage struct{
	// given a subject name, get a list of reports
	subject_to_reports map[string][]*pb.Report
	// given an observer name, get a list of reports
	observer_to_reports map[string][]*pb.Report

	// observers
	observers []Observer

	// Only allow one write at a time
	lock sync.Mutex
}

func NewStorage() Storage{
	return Storage{subject_to_reports:make(map[string][]*pb.Report),observer_to_reports:make(map[string][]*pb.Report),observers:make([]Observer,0)}
}

// register an observer
func (self *Storage) Register(n Observer) (uint64, error){
	self.lock.Lock()
	defer self.lock.Unlock()
	for i, o := range self.observers{
		if o.Module == n.Module && o.Name == n.Name{
			return uint64(i), nil
		}
	}
	self.observers = append(self.observers, n)
	return uint64(len(self.observers) - 1), nil
}

// check if the handle is valid
func (self *Storage) ValidHandle(handle uint64) (bool){
	if handle >= 0 && handle < uint64(len(self.observers)){
		return true
	}
	return false
}

// submit a report to in-mem db
func (self *Storage) SubmitReport(r *pb.Report){
	self.lock.Lock()
	defer self.lock.Unlock()
	self.observer_to_reports[r.Observer] = append(self.observer_to_reports[r.Observer], r)
	self.subject_to_reports[r.Subject] = append(self.subject_to_reports[r.Subject], r)
	//fmt.Println(self.subject_to_reports)
}

// get the latest report on a subject
func (self *Storage) GetLatestReport(subject string) *pb.Report{
	v, found := self.subject_to_reports[subject]
	if found == false{
		return nil
	} else {
		return v[len(v)-1]
	}
}

// get the view on a subject from an observer
func (self *Storage) GetView(subject string, observer string) *pb.View{
	all_reports, found := self.observer_to_reports[observer]
	if found == false{
		return nil
	} else {
		relavant_observations := make([]*pb.Observation,0)
		for _,r := range all_reports{
			if r.Subject == subject && r.Observer == observer{
				relavant_observations = append(relavant_observations,r.Observation)
			}
		}
		if len(relavant_observations) == 0{
			return nil
		} else {
			return &pb.View{Observer:observer,Subject:subject,Observations:relavant_observations}
		}
	}
}

// get the Panorama (full view) on a subject
func (self *Storage) GetPanorama(subject string) *pb.Panorama{
	all_reports, found := self.subject_to_reports[subject]
	if found == false || len(all_reports) == 0 {
		return nil
	} else {
		observations_by_view := make(map[string][]*pb.Observation)
		for _,r := range all_reports{
			observations_by_view[r.Observer] = append(observations_by_view[r.Observer],r.Observation)
		}

		// construct all views
		views := make(map[string]*pb.View,0)
		for observer, observations := range observations_by_view{
			new_view := &pb.View{Observer:observer,Subject:subject,Observations:observations}
			views[observer] = new_view
		}
		return &pb.Panorama{Subject:subject,Views:views}
	}
}

func (self *Storage) DumpPanorama() map[string]*pb.Panorama{
	paranoramas := make(map[string]*pb.Panorama)
	for subject, _ := range self.subject_to_reports{
		paranoramas[subject] = self.GetPanorama(subject)
	}
	return paranoramas
}

// use outer function instead
// func (self *Storage) DumpInference() map[string]*pb.Inference{
// 	paranoramas = make(map[string]*pb.Panorama)
// 	for subject, _ := range self.subject_to_reports{
// 		paranoramas[subject] = self.Inference(subject)
// 	}
// 	return paranoramas
// }

func (self *Storage) DumpSubjects() map[string]*ts.Timestamp{
	subjects := make(map[string]*ts.Timestamp)
	for subject, reports := range self.subject_to_reports{
		subjects[subject] = reports[len(reports)-1].Observation.Ts
	}
	return subjects
}

func (self *Storage) IsObserving(subject string) bool{
	for s,_ := range self.subject_to_reports{
		if s == subject{
			return true
		}
	}
	return false
}