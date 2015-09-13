package main

import (
	"errors"
	"flag"
	log "github.com/Sirupsen/logrus"
	"net"
	"regexp"
)

var regLog *regexp.Regexp

type SteamLog struct {
	Username string
	UID      string
	Wanid    string
	Team     string
	Message  string
}

func init() {
	var err error
	regLog, err = regexp.Compile(`\"(.*?)<(\d{1,5})><\[?(.*?:\d:\d*?)\]?><(.*?)>\"(.*)`)

	if err != nil {
		log.Error("An error occured while compiling regular expressions")
	}
}

func parseLogLine(s string) (SteamLog, error) {
	logMatch := regLog.FindStringSubmatch(s)

	if len(logMatch) > 0 {
		if !T.Check(logMatch[4]) {
			T.Set(logMatch[4], NewTeam())
		}
		return SteamLog{
			Username: logMatch[1],
			UID:      logMatch[2],
			Wanid:    logMatch[3],
			Team:     logMatch[4],
			Message:  logMatch[5],
		}, nil
	} else {
		return SteamLog{}, errors.New("Unable to match log")
	}
}

var T Teams

func main() {

	// SETUP
	log.SetLevel(log.InfoLevel)

	// COMMAND LINE OPTIONS
	optRconaddr := flag.String("rconserver", "127.0.0.1:27015", "server:port of the RCON server to use")
	optRconpass := flag.String("rconpassword", "", "Password of the rcon server")
	optPort := flag.String("loglistener", "127.0.0.1:27500", "Where to listen for server logs")
	optChallengeFile := flag.String("challenges", "", "JSON file of challenge information.")
	flag.Parse()
	rconaddr := *optRconaddr
	rconpass := *optRconpass
	chalfile := *optChallengeFile

	T = NewTeams()

	if rconaddr == "" || rconpass == "" {
		println("ERROR: You did not specify the an address or password for the RCON server")
		flag.PrintDefaults()
		return
	}

	if chalfile == "" {
		println("ERROR: You did not specify a challenge file to load.")
		flag.PrintDefaults()
		return
	}

	// LOAD UP CHALLENGE DATA
	loadChallenges(*optChallengeFile)

	// START ALL OF OUR PROCESSING CHANNELS FOR RESOURCES
	go RunEventChannel()
	go RunRconChannel(rconaddr, rconpass)
	go RunGameChannel()
	go runAPI()

	// LISTEN FOR UDP
	udp, err := net.ResolveUDPAddr("udp", *optPort)
	if err != nil {
		log.Error("Log listener does not appear to be a valid address")
		return
	}
	conn, err := net.ListenUDP("udp", udp)
	if err != nil {
		log.Error("Unable to stard up UDP listener" + err.Error())
		return
	}
	log.WithField("port", *optPort).Info("Listening for UDP connections")

	b := make([]byte, 1500)
	for {
		n, _, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Error("UDP read error: " + err.Error())
			continue
		}

		msg := string(b[:n])
		log.WithField("msg", msg).Debug("Received UDP data")
		l, err := parseLogLine(msg)
		if err != nil {
			continue
		}

		/*LogChan <- Event{
			Message: msg,
		}
		*/
		GameChan <- l
	}
}
