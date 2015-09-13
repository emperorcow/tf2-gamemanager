package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

var EventChan chan Event

type Event struct {
	Message string
}

func init() {
	EventChan = make(chan Event)
}

func RunEventChannel() {
	f, err := os.OpenFile("gamemanager.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Error("Unable to open log file.")
	}
	log.Info("Starting up log channel")
	defer f.Close()

	for {
		l := <-EventChan
		_, err := f.WriteString(l.Message)
		if err != nil {
			log.WithField("msg", l.Message).Error("Unable to write to log file.")
		}
	}
}
