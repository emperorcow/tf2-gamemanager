package main

import (
	log "github.com/Sirupsen/logrus"
	"regexp"
	"strconv"
)

var GameChan chan SteamLog
var regJoinTeam *regexp.Regexp
var regChangeRole *regexp.Regexp
var regKill *regexp.Regexp
var regCapture *regexp.Regexp

type GameRegex struct {
	Name        string
	Expr        string
	Comp        *regexp.Regexp
	CreditValue int
	ScoreValue  int
}

func NewGameRegex(name string, e string, credit int, score int) (GameRegex, error) {
	r, err := regexp.Compile(e)

	if err != nil {
		return GameRegex{}, err
	}

	return GameRegex{
		Name:        name,
		Expr:        e,
		Comp:        r,
		CreditValue: credit,
		ScoreValue:  score,
	}, nil
}

func AppendGameRegex(n string, e string, credit int, score int) {
	tmp, err := NewGameRegex(n, e, credit, score)
	if err != nil {
		log.WithField("pattern", e).Error("An error occured compiling regular expression.")
		return
	}

	LogMessageMatches = append(LogMessageMatches, tmp)
}

var LogMessageMatches []GameRegex

func init() {
	GameChan = make(chan SteamLog)

	var err error
	regJoinTeam, err = regexp.Compile(`joined team \"(.*?)\"`)
	regChangeRole, err = regexp.Compile(`changed role to \"(.*?)\"`)

	if err != nil {
		log.Error("An error occured while compiling regular expressions")
	}

	AppendGameRegex("Player kill", `killed "(.*?)" with "(.*)"`, 1, 50)
	AppendGameRegex("Capture flag", `triggered "flagevent" \(event "captured"\)`, 5, 200)
}

func TeamEventAction(teamname string, username string, score int, credit int) {
	T.AddCredit(teamname, credit)
	T.AddScore(teamname, username, score)
}

func RunGameChannel() {
	for {
		l := <-GameChan

		log.WithField("msg", l.Message).Debug("Game channel got message")

		m := regJoinTeam.FindStringSubmatch(l.Message)
		if len(m) > 0 {
			uid, _ := strconv.ParseInt(l.UID, 10, 64)
			u := User{
				Wanid: l.Wanid,
				Uid:   uid,
			}
			log.WithFields(log.Fields{
				"team":  l.Team,
				"match": m[1],
				"user":  l.Username,
			}).Info("Processing user join team.")

			if !T.Check(m[1]) {
				T.Set(m[1], NewTeam())
			}
			T.AddUser(m[1], l.Username, u)
			continue
		}
		m = regChangeRole.FindStringSubmatch(l.Message)
		if len(m) > 0 {
			log.WithFields(log.Fields{
				"team": l.Team,
				"user": l.Username,
				"role": m[1],
			}).Info("Processing role change")
			T.SetUserRole(l.Team, l.Username, m[1])
			continue
		}
		for _, reg := range LogMessageMatches {
			m = reg.Comp.FindStringSubmatch(l.Message)
			if m != nil {
				log.WithField("name", reg.Name).Info("Matched log event.")
				TeamEventAction(l.Team, l.Username, reg.ScoreValue, reg.CreditValue)
				continue
			}
		}

		log.WithField("m", l.Message).Warn("Unable to match a log event.")
	}
}
