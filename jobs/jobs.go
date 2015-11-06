package jobs

import (
	"fmt"
	"log"
	"time"

	"github.com/mpgerlek/piazza-simulator/message"
)

type JobId int64

var currentId JobId = 1

type JobStatus int

const (
	StatusSubmitted JobStatus = 1
	StatusDispatched JobStatus = 2
	StatusRunning JobStatus = 3
	StatusCompleted JobStatus = 4
	StatusFailed JobStatus = 5
)

type Job struct {
	id        JobId
	status    JobStatus
	messageId MessageId
}

type JobTable struct {
	table map[JobId]Job
}

func NewJobTable() *JobTable {
	var t JobTable
	t.table = make(map[JobId]Job)
	return &t
}

func (t *JobTable) Add(job *Job) {
	if _, ok := t.table[job.id]; ok {
		panic("yow")
	}
	t.table[job.id] = *job
}

func (t *JobTable) Dump() {
	for k, v := range t.table {
		log.Printf("key=%v  value=%v\n", k, v)
	}
}

func NewJob(jobType JobType) *Job {
	var job = Job{id: currentId}
	currentId++

	job.jobType = jobType

	return &job
}

func (job Job) String() string {
	return fmt.Sprintf("{id:%v jobType:%v}", job.id, job.jobType)
}
