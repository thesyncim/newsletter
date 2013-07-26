package main

import (
	"github.com/coocood/jas"
	"labix.org/v2/mgo/bson"
	"log"
	"sync"
	"time"
)

const TickTime = 5

const MaxJobs = 20

var once sync.Once

type Newsletter struct {
	counter int

	//currentJobs
	//completedJobs
	//pendingJobs
}

func (n *Newsletter) PostAdd(ctx *jas.Context) {

	job := new(Job)
	ctx.Unmarshal(job)
	job.Id = bson.NewObjectId()
	job.Status = JOBAdded
	log.Println(JOBRunning)

	err := Col_newsletter_jobs.Insert(job)
	if err != nil {
		log.Println(err)
	}

	Cjob := make(chan Job)

	HandleJobs(Cjob, *job)

	ctx.Data = bson.M{"id": job.Id}

}

func (n *Newsletter) Get(ctx *jas.Context) {

}

/**
*
*
 */
func AddJob(Cjob chan Job, job Job) {

	Cjob <- job

}

/**
*
*
 */
func DispatchJob(Cjob chan Job) {

	select {

	case job := <-Cjob:

		jobTimeUnix := job.ProcessDate.UTC().Truncate(time.Second).Unix()
		//todo fix the need to add +1 Hour to the current time maybe timezone issue
		nowUnix := time.Now().UTC().Add(1 * time.Hour).Truncate(time.Second).Unix()

		/**
		* if job shedule time is higher than Current Time we need to Start the
		* job immediately otherwise we add it to the PendingJobs queue
		 */
		log.Println(jobTimeUnix, nowUnix)

		if jobTimeUnix <= nowUnix {

			go StartJob(job)

		} else {
			log.Println("reach enqueue")

			PendingJobs.mutex.Lock()
			PendingJobs.Jobs[job.Id] = job
			PendingJobs.mutex.Unlock()

			//start the Pending jobs handler only once
			go func() {
				once.Do(HandlePendingJobs)
			}()

			log.Println(PendingJobs.Jobs)
		}
	}
}

/**
* this function is executed only once after we receive the first pending job
*
 */
func HandlePendingJobs() {

	c := time.Tick(TickTime * time.Second)

	for now := range c {
		PendingJobs.mutex.Lock()
		for key, job := range PendingJobs.Jobs {

			jobTimeUnix := job.ProcessDate.UTC().Truncate(time.Second).Unix()
			nowUnix := now.UTC().Add(1 * time.Hour).Truncate(time.Second).Unix()
			log.Println(nowUnix, jobTimeUnix)

			if nowUnix > jobTimeUnix {
				go StartJob(job)
				delete(PendingJobs.Jobs, key)
				log.Println("key ", key, " removed")
			}

		}
		PendingJobs.mutex.Unlock()
	}
}

/**
*
*
 */
func StartJob(job Job) {

	var record NewsletterRecord

	newslettercontent := Col_newsletter.FindId(bson.ObjectIdHex(job.NewsletterId))
	err := newslettercontent.One(&record)
	if err != nil {
		log.Println(job.NewsletterId, err)
	}

	var listrecord []NewsletterListRecord

	listrecord = make([]NewsletterListRecord, len(job.EmailList))

	for index, id := range job.EmailList {
		Col_newsletter_lists.FindId(bson.ObjectIdHex(id)).One(&listrecord[index])
	}

	sendEmail(listrecord, record)

}

func HandleJobs(cjob chan Job, job Job) {

	go AddJob(cjob, job)

	go DispatchJob(cjob)

}
