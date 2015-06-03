package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
)

var Challenges []Challenge

type Challenge struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Difficulty  string `json:"difficulty"`
	Description string `json:"description"`
	Value       int    `json:"value"`
}

func loadChallenges(filepath string) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.WithField("path", filepath).Error("Unable to open challenge file.")
	}

	Challenges = make([]Challenge, 0)
	json.Unmarshal(file, &Challenges)

	log.WithField("challenges", Challenges).Debug("Challenges loaded.")
}

func GetChallenge(id string) (Challenge, bool) {
	for _, chal := range Challenges {
		if chal.ID == id {
			return chal, true
		}
	}
	return Challenge{}, false
}
