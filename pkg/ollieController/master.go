package ollieController

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
)

// OllieCommand to describe command to be executed on Ollie
// Fields:
// 	Command: 	Command to be executed
//	Direction: 	Direction to move ollie for commands like Roll, SetHeading, etc
//		   	set 0 if command doesn't support direction parameter
//	Speed: 		Speed with which ollie should move for commands like SetMotorValues, Roll, etc
//			set 0 if command doesn't support speed parameter
// 	Duration:	Time in Milliseconds for which command will be executed
//			set 0 to run command for forever
type OllieCommand struct {
	Command   string
	Direction uint16
	Speed     uint8
	Duration  uint16
}

var (
	ch chan OllieCommand
)

func SendCommand(cmd OllieCommand) {
	ch <- cmd
}

func InitController(ollies []string) {
	master := gobot.NewMaster()
	api.NewAPI(master).Start()
	ch = make(chan OllieCommand, 2)
	for _, port := range ollies {
		master.AddRobot(NewOllieBot(port))
	}
	go master.Start()
}
