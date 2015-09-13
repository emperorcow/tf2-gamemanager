package main

import ()

type TeamAction struct {
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Cost        int    `json:"cost"`
	Description string `json:"description"`
	Cmd         string `json:"cmd"`
	TargetSelf  bool   `json:"targetself"`
	TargetRand  bool   `json:"targetrand"`
}

var Actions = []TeamAction{
	TeamAction{
		"Low Gravity",
		"fighter-jet",
		5,
		"Your whole team can fly for 60 seconds",
		"sm_forcertd %s 5",
		true,
		false,
	},
	TeamAction{
		"Invulnerable",
		"star",
		20,
		"Your team become uber for 15 seconds, unable to be killed.",
		"sm_forcertd %s 0",
		true,
		false,
	},
	TeamAction{
		"Crits",
		"bolt",
		9,
		"For 30 seconds, all of your teams shots become critical hits (double damage).",
		"sm_forcertd %s 9",
		true,
		false,
	},
	TeamAction{
		"Infinite Ammo",
		"crosshairs",
		11,
		"Your entire team has unlimited ammunition for 30 seconds.",
		"sm_forcertd %s 10",
		true,
		false,
	},
	TeamAction{
		"The Bomb!",
		"bomb",
		8,
		"A random member of the opposite team becomes 'The Bomb' and will explode in 10 seconds.  Better stay away!",
		"sm_forcertd #%s 18",
		false,
		true,
	},
	TeamAction{
		"Drug Em",
		"flask",
		4,
		"Drug the other team, sending them into a crazy hallucinogenic frenzy, if they can find you for 15 seconds.",
		"sm_forcertd %s 21",
		false,
		false,
	},
	TeamAction{
		"Earthquake!",
		"retweet",
		4,
		"Shake the other teams screen for 15 seconds",
		"sm_forcertd %s 27",
		false,
		false,
	},
	TeamAction{
		"Instakill",
		"heartbeat",
		50,
		"A random member of your team becomes death incarnate, every shot will instantly kill whomever or whatever it hits.",
		"sm_forcertd #%s 32",
		true,
		true,
	},
	TeamAction{
		"Honey, I Shrunk the Consultant",
		"level-down",
		6,
		"A random member of your team becomes tiny, maybe they can sneak in behind the other team?",
		"sm_forcertd #%s 34",
		true,
		true,
	},
	TeamAction{
		"Spin the Wheel",
		"question-circle",
		2,
		"Hmm... something seems broken here.  Every time I click this button something random happens to someone on my team.  Well... good luck!",
		"sm_forcertd #%s",
		true,
		true,
	},
	TeamAction{
		"Burn Em!",
		"fire",
		8,
		"Set the other team on fire for 15 seconds!  Seems straightforward.",
		"sm_burn %s 15",
		false,
		false,
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
