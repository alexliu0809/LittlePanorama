package store

import (
	"fmt"
	"sync"

	dd "panorama/decision"
	dt "panorama/types"
	du "panorama/util"

	pb "panorama/build/gen"
)

const (
	itag = "inference"
)

type InferMap map[string]*pb.Inference

type HealthInferenceStorage struct {
	Results   InferMap
	Workbooks map[string]InferMap
	ReportCh  chan *pb.Report
	SubjectCh chan string

	raw   dt.HealthStorage
	db    dt.HealthDB
	algo  dd.InferenceAlgo
	mu    *sync.RWMutex
	alive bool
}

func NewHealthInferenceStorage(raw dt.HealthStorage, algo dd.InferenceAlgo) *HealthInferenceStorage {
	storage := &HealthInferenceStorage{
		Results:   make(InferMap),
		Workbooks: make(map[string]InferMap),
		ReportCh:  make(chan *pb.Report, 50),
		SubjectCh: make(chan string, 50),
		raw:       raw,
		algo:      algo,
		mu:        &sync.RWMutex{},
		alive:     true,
	}
	return storage
}

var _ dt.HealthInference = new(HealthInferenceStorage)

func (self *HealthInferenceStorage) InferSubjectAsync(subject string) error {
	// simply sent it to channel and return
	self.SubjectCh <- subject
	return nil
}

func (self *HealthInferenceStorage) InferReportAsync(report *pb.Report) error {
	// simply sent it to channel and return
	self.ReportCh <- report
	return nil
}

func (self *HealthInferenceStorage) InferSubject(subject string) (*pb.Inference, error) {
	pano := self.raw.GetPanorama(subject)
	if pano == nil {
		du.LogD(itag, "empty panorama for %s, reset inference result to empty", subject)
		self.mu.Lock()
		delete(self.Workbooks, subject)
		delete(self.Results, subject)
		self.mu.Unlock()
		return nil, fmt.Errorf("cannot get panorama for %s\n", subject)
	}
	// since we need to re-calculate the inference for the entire subject
	// we should clear the workbook
	workbook := make(InferMap)
	self.mu.Lock()
	self.Workbooks[subject] = workbook
	self.mu.Unlock()
	pano.RLock()
	inference := self.algo.InferPano(pano.Value, workbook)
	pano.RUnlock()
	if inference == nil {
		du.LogD(itag, "empty inference for %s, reset result to empty", subject)
		self.mu.Lock()
		delete(self.Workbooks, subject)
		delete(self.Results, subject)
		self.mu.Unlock()
		return nil, fmt.Errorf("could not compute inference for %s\n", subject)
	}
	// du.LogD(itag, "inference result for %s: %s", subject, dt.ObservationString(inference.Observation))
	self.mu.Lock()
	self.Results[subject] = inference
	self.mu.Unlock()
	return inference, nil
}

func (self *HealthInferenceStorage) InferReport(report *pb.Report) (*pb.Inference, error) {
	// TODO: support incremental inference
	pano := self.raw.GetPanorama(report.Subject)
	if pano == nil {
		du.LogD(itag, "empty panorama for %s, reset inference result to empty", report.Subject)
		self.mu.Lock()
		delete(self.Workbooks, report.Subject)
		delete(self.Results, report.Subject)
		self.mu.Unlock()
		return nil, fmt.Errorf("cannot get panorama for %s\n", report.Subject)
	}
	self.mu.Lock()
	workbook, ok := self.Workbooks[report.Subject]
	if !ok {
		workbook = make(InferMap)
		self.Workbooks[report.Subject] = workbook
	} else {
		// clear the workbook entry for the particular observer
		// so that we just need to re-infer the specific view
		delete(workbook, report.Observer)
	}
	self.mu.Unlock()
	pano.RLock()
	inference := self.algo.InferPano(pano.Value, workbook)
	pano.RUnlock()
	if inference == nil {
		du.LogD(itag, "empty inference for %s, reset result to empty", report.Subject)
		self.mu.Lock()
		delete(self.Workbooks, report.Subject)
		delete(self.Results, report.Subject)
		self.mu.Unlock()
		return nil, fmt.Errorf("could not compute inference for %s\n", report.Subject)
	}
	du.LogD(itag, "inference result for %s: %s", report.Subject, dt.ObservationString(inference.Observation))
	self.mu.Lock()
	self.Results[report.Subject] = inference
	self.mu.Unlock()
	return inference, nil
}

func (self *HealthInferenceStorage) GetInference(subject string) *pb.Inference {
	self.mu.Lock()
	inference, ok := self.Results[subject]
	self.mu.Unlock()
	if !ok {
		return nil
	}
	return inference
}

func (self *HealthInferenceStorage) DumpInference() map[string]*pb.Inference {
	return self.Results
}

func (self *HealthInferenceStorage) SetDB(db dt.HealthDB) {
	self.db = db
}

func (self *HealthInferenceStorage) Start() error {
	go func() {
		for self.alive {
			select {
			case subject := <-self.SubjectCh:
				{
					du.LogD(itag, "perform inference on subject for %s", subject)
					inf, err := self.InferSubject(subject)
					if err != nil {
						du.LogE(itag, "failed to infer for %s", subject)
					} else {
						if self.db != nil && inf != nil {
							go self.db.InsertInference(inf)
						}
					}
				}
			case report := <-self.ReportCh:
				{
					if len(report.Subject) != 0 {
						du.LogD(itag, "received report for %s for inference", report.Subject)
						inf, err := self.InferReport(report)
						if err != nil {
							du.LogE(itag, "failed to infer for %s", report.Subject)
						} else {
							if self.db != nil && inf != nil {
								go self.db.InsertInference(inf)
							}
						}
					}
				}
			}
		}
	}()
	return nil
}

func (self *HealthInferenceStorage) Stop() error {
	self.alive = false
	var report pb.Report
	select {
	case self.ReportCh <- &report:
		du.LogI(itag, "send empty report to stop the service")
	default:
	}
	return nil
}
