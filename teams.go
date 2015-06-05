package main

import (
	log "github.com/Sirupsen/logrus"
	"sync"
)

type Team struct {
	Pass       string          `json:"pass"`
	Users      Users           `json:"users"`
	Score      int             `json:"score"`
	Credits    int             `json:"credits"`
	Challenges map[string]bool `json:"challenges"`
}

type Teams struct {
	data map[string]Team
	sync.RWMutex
}

func NewTeams() Teams {
	return Teams{
		data: make(map[string]Team),
	}
}

func (t *Teams) Set(name string, info Team) {
	log.WithField("name", name).Debug("Set team")
	t.Lock()
	t.data[name] = info
	t.Unlock()
}

func (t *Teams) Get(name string) (Team, bool) {
	log.WithField("name", name).Debug("Get team")
	t.RLock()
	tmp, ok := t.data[name]
	t.RUnlock()
	return tmp, ok
}

type Pair struct {
	Key string
	Val Team
}

func (t *Teams) Iterate() <-chan Pair {
	ch := make(chan Pair)
	go func() {
		t.RLock()
		for k, v := range t.data {
			ch <- Pair{k, v}
		}
		t.RUnlock()
		close(ch)
	}()
	return ch
}

func (t *Teams) Check(name string) bool {
	t.RLock()
	_, ok := t.data[name]
	t.RUnlock()
	log.WithFields(log.Fields{
		"name":   name,
		"result": ok,
	}).Debug("Check team")
	return ok
}

func (t *Teams) Delete(name string) {
	log.WithField("name", name).Debug("Delete team")
	t.Lock()
	delete(t.data, name)
	t.Unlock()
}

func (t *Teams) SetChallenge(teamname string, challengeid string, status bool) {
	t.Lock()
	t.data[teamname].Challenges[challengeid] = status
	t.Unlock()
}

func (t *Teams) SetPassword(name string, pass string) {
	log.WithFields(log.Fields{
		"name": name,
		"pass": pass,
	}).Debug("Set team password")
	t.Lock()
	tmp := t.data[name]
	tmp.Pass = pass
	t.data[name] = tmp
	t.Unlock()
}

func (t *Teams) AddScore(team string, user string, amt int) {
	log.WithFields(log.Fields{
		"team": team,
		"amt":  amt,
	}).Info("Added score to team")
	t.Lock()
	tmp := t.data[team]
	tmp.Score += amt
	tmp.Users.AddScore(user, amt)
	t.data[team] = tmp
	t.Unlock()
}

func (t *Teams) AddCredit(team string, amt int) {
	log.WithFields(log.Fields{
		"team": team,
		"amt":  amt,
	}).Info("Added credit to team")
	t.Lock()
	tmp := t.data[team]
	tmp.Credits += amt
	t.data[team] = tmp
	t.Unlock()
}

func (t *Teams) AddUser(teamname string, username string, u User) {
	t.Lock()
	tmp := t.data[teamname]
	tmp.Users.Add(username, u)
	t.Unlock()
}

func (t *Teams) SetUserRole(teamname string, username string, role string) {
	t.Lock()
	tmp := t.data[teamname]
	tmp.Users.SetRole(username, role)
	t.Unlock()
}

func (t *Teams) GetUser(teamname string, username string) (User, bool) {
	t.RLock()
	tmp := t.data[teamname]
	u, ok := tmp.Users.Get(username)
	t.RUnlock()
	return u, ok
}

func (t *Teams) DelUser(teamname string, username string) {
	t.Lock()
	tmp := t.data[teamname]
	tmp.Users.Delete(username)
	t.Unlock()
}
