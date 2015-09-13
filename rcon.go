package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kidoman/go-steam"
)

var RconChan chan RconData

type RconData struct {
	Cmd  string
	Resp chan string
}

func init() {
	RconChan = make(chan RconData, 5)
}

func RunRconChannel(addr string, pass string) {
	opts := &steam.ConnectOptions{
		RCONPassword: pass,
	}

	for {
		r := <-RconChan

		rcon, err := steam.Connect(addr, opts)
		defer rcon.Close()

		if err != nil {
			log.WithField("error", err.Error()).Error("Unable to connect to RCON server.")
			continue
		}

		resp, err := rcon.Send(r.Cmd)
		if err != nil {
			log.WithField("cmd", r.Cmd).Error("Unable to execute RCON command: " + err.Error())
			continue
		}

		log.WithFields(log.Fields{
			"cmd":  r.Cmd,
			"resp": resp,
		}).Info("RCON command successful.")
		r.Resp <- resp
	}
}
