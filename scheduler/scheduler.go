package scheduler

import (
	"encoding/json"

	"github.com/colinyl/cron"
)

type Scheduler struct {
	c *cron.Cron
}

var (
	Schd *Scheduler = &Scheduler{c: cron.New()}
)
func AddTask(trigger string,task *TaskDetail){
    Schd.c.AddJob(trigger,task)
}

func AddJob(job *JobDetail) {
	Schd.c.AddJob(job.config.Trigger, job)
}

func Start() {
	Schd.c.Start()
}

func Stop() {
	Schd.c.Stop()
	Schd = &Scheduler{c: cron.New()}
}

func (s *Scheduler) GetSnap() string {
	snaps := make(map[string]*JobSnap)
	entries := s.c.Entries()
	for _, v := range entries {
		job := v.Job.(*JobDetail)
		snaps[job.config.Name] = &JobSnap{Name: job.config.Name,
			Prev: v.Prev, Next: v.Next}

	}
	buffer, _ := json.Marshal(snaps)
	return string(buffer)

}
