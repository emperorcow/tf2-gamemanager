package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

const (
	// Integer Status Codes
	HTTP_CODE_OK           = 200
	HTTP_CODE_CREATED      = 201
	HTTP_CODE_NOCONTENT    = 204
	HTTP_CODE_NOTMODIFIED  = 304
	HTTP_CODE_BADREQ       = 400
	HTTP_CODE_UNAUTHORIZED = 401
	HTTP_CODE_FORBIDDEN    = 403
	HTTP_CODE_NOTFOUND     = 404
	HTTP_CODE_CONFLICT     = 409
	HTTP_CODE_ERROR        = 500

	// Text Status Codes
	HTTP_CODE_OK_T           = "OK"
	HTTP_CODE_CREATED_T      = "Created"
	HTTP_CODE_NOCONTENT_T    = "No Content"
	HTTP_CODE_NOTMODIFIED_T  = "Not Modified"
	HTTP_CODE_BADREQ_T       = "Bad Request"
	HTTP_CODE_UNAUTHORIZED_T = "Unauthorized"
	HTTP_CODE_FORBIDDEN_T    = "Forbidden"
	HTTP_CODE_NOTFOUND_T     = "Not Found"
	HTTP_CODE_CONFLICT_T     = "Conflict"
	HTTP_CODE_ERROR_T        = "Internal Server Error"
)

var subnetAdmin *net.IPNet
var subnetTeam1 *net.IPNet
var subnetTeam2 *net.IPNet

func init() {
	_, subnetAdmin, _ = net.ParseCIDR("192.168.0.0/24")
	_, subnetTeam1, _ = net.ParseCIDR("192.168.1.0/24")
	_, subnetTeam2, _ = net.ParseCIDR("192.168.2.0/24")
}

func SubnetMatcher(addr string, n *net.IPNet) bool {
	ip := net.ParseIP(addr)
	iplog := log.WithField("ip", addr)

	if ip == nil {
		iplog.Warn("Unable to parse remote requesting ip")
		return false
	}
	if ip.To4() == nil {
		iplog.Warn("Address is not an IPv4 address")
		return false
	}

	if n.Contains(ip) {
		return true
	} else {
		iplog.Warn("IP is not within allowed range")
		return false
	}
}

func SubnetMatchAdmin(r *http.Request, rm *mux.RouteMatch) bool {
	return SubnetMatcher(r.RemoteAddr, subnetAdmin)
}
func SubnetMatchTeam1(r *http.Request, rm *mux.RouteMatch) bool {
	return SubnetMatcher(r.RemoteAddr, subnetTeam1)
}
func SubnetMatchTeam2(r *http.Request, rm *mux.RouteMatch) bool {
	return SubnetMatcher(r.RemoteAddr, subnetTeam2)
}

func runAPI() {
	r := mux.NewRouter()

	r.HandleFunc("/api/teams", apiTeamsQuery).Methods("GET")
	r.HandleFunc("/api/teams/{team}", apiTeamsGET).Methods("GET")
	r.HandleFunc("/api/teams", apiTeamsPUT).Methods("PUT")
	r.HandleFunc("/api/teams/{team}", apiTeamsDELETE).Methods("DELETE")

	r.HandleFunc("/api/actions", apiActionsGet).Methods("GET")
	r.HandleFunc("/api/actions/{name}", apiActionsSet).Methods("GET")

	r.HandleFunc("/api/challenges", apiChallengesGet).Methods("GET")
	r.HandleFunc("/api/challenges", apiChallengesSet).Methods("POST")

	r.HandleFunc("/api/rcon", apiRconPost).Methods("POST")

	r.PathPrefix("/").MatcherFunc(SubnetMatchAdmin).Handler(http.FileServer(http.Dir("./static/admin/")))
	r.PathPrefix("/").MatcherFunc(SubnetMatchTeam1).Handler(http.FileServer(http.Dir("./static/team1/")))
	//	r.PathPrefix("/").MatcherFunc(SubnetMatchTeam2).Handler(http.FileServer(http.Dir("./static/team2/")))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/team2/")))

	http.ListenAndServe(":80", r)
}

type apiTeams struct {
	Name string `json:"name"`
	Info Team   `json:"info"`
}

func apiTeamsQuery(w http.ResponseWriter, r *http.Request) {
	var data []apiTeams
	resp := json.NewEncoder(w)

	for p := range T.Iterate() {
		log.WithFields(log.Fields{
			"key":  p.Key,
			"data": p.Val,
		}).Debug("Added team to response")

		if p.Key != "" && p.Key != "Unassigned" && p.Key != "Spectator" {
			data = append(data, apiTeams{p.Key, p.Val})
		}
	}

	w.WriteHeader(HTTP_CODE_OK)
	resp.Encode(data)
	log.Info("Provided a team listing to API")
}
func apiTeamsGET(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["team"]
	resp := json.NewEncoder(w)

	team, ok := T.Get(name)
	if !ok {
		w.WriteHeader(HTTP_CODE_NOTFOUND)
		log.WithField("name", name).Error("Unable to find team requested")
		return
	}

	var tmp apiTeams
	tmp.Name = name
	tmp.Info = team

	w.WriteHeader(HTTP_CODE_OK)
	resp.Encode(tmp)
	log.WithField("name", name).Debug("Provided team info to API")
}
func apiTeamsPUT(w http.ResponseWriter, r *http.Request) {
	var recData apiTeams
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1000000))
	if err != nil {
		w.WriteHeader(HTTP_CODE_FORBIDDEN)
		log.Error("The data received to create a team was too large.")
		return
	}
	if err := json.Unmarshal(body, &recData); err != nil {
		w.WriteHeader(HTTP_CODE_BADREQ)
		log.Error("400: Unable to process data sent in put body for team create.")
		return
	}

	T.Set(recData.Name, recData.Info)
	w.WriteHeader(HTTP_CODE_OK)
}
func apiTeamsPOST(w http.ResponseWriter, r *http.Request) {
	var recData apiTeams
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1000000))
	if err != nil {
		w.WriteHeader(HTTP_CODE_FORBIDDEN)
		log.Error("The data received to update a team was too large.")
		return
	}
	if err := json.Unmarshal(body, &recData); err != nil {
		w.WriteHeader(HTTP_CODE_BADREQ)
		log.Error("400: Unable to process data sent in post body to update team.")
		return
	}

	T.Set(recData.Name, recData.Info)
	w.WriteHeader(HTTP_CODE_OK)
}
func apiTeamsDELETE(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["team"]
	T.Delete(name)
	w.WriteHeader(HTTP_CODE_OK)
}

func apiActionsGet(w http.ResponseWriter, r *http.Request) {
	log.Info("Sent available actions for a team to API.")
	resp := json.NewEncoder(w)
	w.WriteHeader(HTTP_CODE_OK)
	resp.Encode(Actions)
}

func apiActionsSet(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	l := log.WithFields(log.Fields{
		"action": name,
	})

	act, ok := getAction(name)
	if !ok {
		w.WriteHeader(HTTP_CODE_ERROR)
		fmt.Fprintf(w, "Cannot find that action")
		l.Error("Unable to run action from team.")
		return
	}

	var target string
	if !act.TargetRand {
		if SubnetMatcher(r.RemoteAddr, subnetTeam1) {
			if act.TargetSelf {
				target = "@red"
			} else {
				target = "@blue"
			}
		} else if SubnetMatcher(r.RemoteAddr, subnetTeam2) {
			if act.TargetSelf {
				target = "@blue"
			} else {
				target = "@red"
			}
		} else {
			w.WriteHeader(HTTP_CODE_BADREQ)
			fmt.Fprintf(w, "Request did not come from a proper subnet")
			l.Error("Request did not come from a team subnet.")
			return
		}
	} else {
		if SubnetMatcher(r.RemoteAddr, subnetTeam1) {
			if act.TargetSelf {
				target, _ = T.GetRandomUser("Red")
			} else {
				target, _ = T.GetRandomUser("Blue")
			}
		} else if SubnetMatcher(r.RemoteAddr, subnetTeam2) {
			if act.TargetSelf {
				target, _ = T.GetRandomUser("Blue")
			} else {
				target, _ = T.GetRandomUser("Red")
			}
		} else {
			w.WriteHeader(HTTP_CODE_BADREQ)
			fmt.Fprintf(w, "Request did not come from a proper subnet")
			l.Error("Request did not come from a team subnet.")
			return
		}
	}

	respChan := make(chan string, 1)
	cmd := fmt.Sprintf(act.Cmd, target)
	RconChan <- RconData{cmd, respChan}

	out := <-respChan

	l.WithField("output", out).Info("Ran action for team against target.")
	w.WriteHeader(HTTP_CODE_OK)
	fmt.Fprintf(w, out)
}

func apiChallengesGet(w http.ResponseWriter, r *http.Request) {
	log.Info("Sent challenge information to API.")
	resp := json.NewEncoder(w)
	w.WriteHeader(HTTP_CODE_OK)
	resp.Encode(Challenges)
}

type apiChallenge struct {
	ID     string `json:"id"`
	Team   string `json:"team"`
	Status bool   `json:"status"`
}

func apiChallengesSet(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 100000))
	if err != nil {
		log.Error("There was an gathering data to update challenge status from POST data: " + err.Error())
	}

	var recData apiChallenge
	if err := json.Unmarshal(body, &recData); err != nil {
		w.WriteHeader(HTTP_CODE_BADREQ)
		log.Error("400: Unable to process data sent in put body for challenge set.")
		return
	}

	chal, ok := GetChallenge(recData.ID)
	if !ok {
		w.WriteHeader(HTTP_CODE_ERROR)
		fmt.Fprintf(w, "Cannot find that team or challenge")
		return
	}

	log.WithFields(log.Fields{
		"team":      recData.Team,
		"challenge": recData.ID,
		"status":    recData.Status,
	}).Info("Setting challenge status for team.")

	T.SetChallenge(recData.Team, recData.ID, recData.Status)
	if recData.Status {
		T.AddCredit(recData.Team, chal.Value)
		T.AddScore(recData.Team, "Puzzlers", chal.Value*10)
	} else {
		T.AddCredit(recData.Team, chal.Value*-1)
		T.AddScore(recData.Team, "Puzzlers", chal.Value*-10)
	}

	w.WriteHeader(HTTP_CODE_OK)
}

func apiRconPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 10000))
	if err != nil {
		log.Error("There was an error parsing an RCON command " + err.Error())
	}

	respChan := make(chan string, 1)
	cmd := string(body[:])

	log.WithField("body", cmd).Debug("Gathered command from API")
	RconChan <- RconData{cmd, respChan}

	out := <-respChan

	log.WithField("out", out).Debug("Command output received")

	w.WriteHeader(HTTP_CODE_OK)
	fmt.Fprintf(w, out)
}
