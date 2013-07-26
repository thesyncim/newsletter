package main

import (
	"labix.org/v2/mgo"
	"log"
)

var (
	Col_newsletter_jobs  *mgo.Collection
	Col_newsletter       *mgo.Collection
	Col_newsletter_lists *mgo.Collection
	DBSession            *mgo.Session
)

func init_DB() {
	DBSession, err := mgo.Dial("server1.thesyncim.com")

	if err != nil {
		log.Println(err)
	}

	db := DBSession.DB("theatrix")

	Col_newsletter_jobs = db.C("newsletter_jobs")
	Col_newsletter_lists = db.C("newsletter_lists")
	Col_newsletter = db.C("newsletters")

	if err != nil {
		panic(err)
	}
}
