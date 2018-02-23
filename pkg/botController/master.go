package botController

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
)

// Command to describe command to be executed on bot
// Fields:
// 	Command: 	Command to be executed
//	Direction: 	Direction to move ollie for commands like Roll, SetHeading, etc
//		   	set 0 if command doesn't support direction parameter
//	Speed: 		Speed with which ollie should move for commands like SetMotorValues, Roll, etc
//			set 0 if command doesn't support speed parameter
// 	Duration:	Time in Milliseconds for which command will be executed
//			set 0 to run command for forever
type Command struct {
	Command   string
	Direction int16
	Speed     uint8
	Duration  uint16
}

const (
	defaultSpeed    = uint8(255)
	defaultDir      = 0
	defaultInterval = 100
	DefaultDur      = 5000
)

var (
	DefaultRollSpeed uint8
	cmdDirection     int16
	cmdSpeed         uint8
	cmdDuration      uint16
)

func (olcmd Command) String() string {
	return fmt.Sprintf("{command:%v direction:%v speed:%v duration:%v}", olcmd.Command, olcmd.Direction, olcmd.Speed, olcmd.Duration)
}

var (
	ol chan Command
	sp chan Command
)

func SendCommandToOllie(cmd Command) {
	ol <- cmd
}

func SendCommandToSphero(cmd Command) {
	sp <- cmd
}

func checkDuration(iter *uint16, device string) {
	if *iter > 0 {
		*iter--
	}
	if *iter == 0 {
		if device == "ollie" {
			ol <- Command{"stop", 0, 0, 0}
		} else if device == "sphero" {
			sp <- Command{"stop", 0, 0, 0}
		}
	}
}

func InitController(botName string, botIds []string) {
	master := gobot.NewMaster()
	api.NewAPI(master).Start()
	ol = make(chan Command, 2)
	sp = make(chan Command, 2)
	switch botName {
	case "sphero":
		for _, port := range botIds {
			master.AddRobot(NewSpheroBot(port))
		}
	case "ollie":
		for _, port := range botIds {
			master.AddRobot(NewOllieBot(port))
		}
	}
	go master.Start()
}
