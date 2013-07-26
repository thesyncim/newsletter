// structs
package main

import (
	"labix.org/v2/mgo/bson"
	"sync"
	"time"
)

const (
	JOBAdded     = 0
	JOBPending   = 1
	JOBRunning   = 2
	JOBCompleted = 3
	JOBFailed    = 4
)

var PendingJobs = new(PendingJobsStruct)

type NewsletterRecord struct {
	Id          bson.ObjectId "_id,omitempty"
	Author      string
	Content     string
	Description string
	Title       string
}

type NewsletterListRecord struct {
	Id     bson.ObjectId "_id,omitempty"
	name   string
	Emails []string
}

type PendingJobsStruct struct {
	Jobs    map[bson.ObjectId]Job
	Started chan bool
	mutex   sync.RWMutex
}

type Job struct {
	Id           bson.ObjectId "_id,omitempty"
	EmailList    []string
	SmtpAccount  string
	NewsletterId string
	ProcessDate  time.Time // RFC 3339
	Status       int
}

func init() {
	PendingJobs.Started = make(chan bool)
	PendingJobs.Jobs = make(map[bson.ObjectId]Job)
}
