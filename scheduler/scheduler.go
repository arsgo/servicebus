
package scheduler

import (
	"github.com/colinyl/cron"
)

type Scheduler struct {
	c *cron.Cron
}

var (
	scheduler *Scheduler = &Scheduler{c: cron.New()}
)

func AddJob(job *JobDetail) {
	scheduler.c.AddJob(job.config.Trigger, job)
}

func Start() {
	scheduler.c.Start()
}

func Stop() {
	scheduler.c.Stop()
}
func GetSnap() map[string]*JobSnap {
	snaps := make(map[string]*JobSnap)
	entries := scheduler.c.Entries()
	for _, v := range entries {
		job := v.Job.(*JobDetail)
		snaps[job.config.Name] = &JobSnap{Name: job.config.Name,
			Prev: v.Prev, Next: v.Next}

	}
	return snaps

}
