package main

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"sync"
)

type User struct {
	Wanid string `json:"wanid"`
	Uid   int64  `json:"uid"`
	Role  string `json:"role"`
	Score int    `json:"score"`
}

type Users struct {
	Data map[string]User `json:"users"`
	sync.RWMutex
}

func NewUsers() Users {
	return Users{
		Data: make(map[string]User),
	}
}

func (u *Users) Add(username string, usr User) {
	log.WithField("name", username).Debug("Added user")
	u.Lock()
	u.Data[username] = usr
	u.Unlock()
}

func (u *Users) Get(username string) (User, bool) {
	log.WithField("name", username).Debug("Get user")
	u.RLock()
	usr, ok := u.Data[username]
	u.RUnlock()
	return usr, ok
}

func (u *Users) GetRandom() (string, bool) {
	log.Debug("Get random user")
	u.RLock()
	users := make([]string, 0, len(u.Data))
	if len(u.Data) == 0 {
		return "", false
	}

	for k := range u.Data {
		users = append(users, k)
	}
	randUsername := users[rand.Intn(len(users))]
	u.RUnlock()

	return randUsername, true
}

func (u *Users) GetAll() map[string]User {
	log.Debug("Get all users")
	tmp := make(map[string]User)

	u.RLock()
	for k, v := range u.Data {
		tmp[k] = v
	}
	u.RUnlock()

	return tmp
}

func (u *Users) Check(username string) bool {
	u.RLock()
	_, ok := u.Data[username]
	u.RUnlock()
	log.WithFields(log.Fields{
		"name":   username,
		"result": ok,
	}).Debug("Check user")
	return ok
}

func (u *Users) Delete(username string) {
	log.WithField("name", username).Debug("Delete user")
	u.Lock()
	delete(u.Data, username)
	u.Unlock()
}

func (u *Users) AddScore(username string, score int) {
	log.WithFields(log.Fields{
		"name":  username,
		"score": score,
	}).Debug("Adding to user score")
	u.Lock()
	tmp := u.Data[username]
	tmp.Score += score
	u.Data[username] = tmp
	u.Unlock()
}

func (u *Users) SetRole(username string, role string) {
	log.WithFields(log.Fields{
		"name": username,
		"role": role,
	}).Debug("Set user role")
	u.Lock()
	tmp := u.Data[username]
	tmp.Role = role
	u.Data[username] = tmp
	u.Unlock()
}
