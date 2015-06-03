package main

import ()

type TeamAction struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Cost        int    `json:"cost"`
	Description string `json:"description"`
	Cmd         string `json:"cmd"`
}

var Actions = []TeamAction{
	TeamAction{
		"Burn",
		"fire",
		50,
		"Set the opposite team on fire for 15 seconds",
		"sm_burn @%s 15",
	},
	TeamAction{
		"Bomb",
		"bomb",
		50,
		"Set one random member on the other team to explode in 10 seconds",
		"sm_bomb @%s 10",
	},
	TeamAction{
		"Invulnerability",
		"star",
		50,
		"You're invulnerable for 30 seconds!",
		"sm_burn @%s 15",
	},
}

func getAction(name string) (TeamAction, bool) {
	for _, act := range Actions {
		if act.Name == name {
			return act, true
		}
	}
	return TeamAction{}, false
}
