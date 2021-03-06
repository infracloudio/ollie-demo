package botController

import (
	log "github.com/sirupsen/logrus"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/sphero"
)

func NewSpheroBot(port string) *gobot.Robot {
	cmdDirection = defaultDir
	cmdSpeed = defaultSpeed
	cmdDuration = DefaultDur
	cmdInterval := defaultInterval * time.Millisecond
	adaptor := sphero.NewAdaptor("/dev/rfcomm0")
	spheroBot := sphero.NewSpheroDriver(adaptor)

	work := func() {
		spheroHead := int16(0)
		spheroBot.Roll(0, 0)
		var ticker *time.Ticker
		for {
			cmd := <-sp
			log.WithFields(log.Fields{
				"device": "sphero",
				"cmd":    cmd,
			}).Info("command received by sphero driver")
			// Stop previous command execution
			if ticker != nil {
				spheroBot.SetRGB(255, 0, 0)
				ticker.Stop()
			}
			cmdDirection = cmd.Direction
			if cmd.Speed != 0 {
				cmdSpeed = cmd.Speed
			} else {
				cmdSpeed = defaultSpeed
			}
			if cmd.Duration > 0 {
				cmdDuration = cmd.Duration
			} else {
				cmdDuration = DefaultDur
			}
			// No of command iterations (1 per defaultInterval ms)
			it := cmdDuration / defaultInterval
			switch cmd.Command {
			case "jump":
				spheroBot.SetRGB(0, 0, 255)
				spheroBot.Roll(255, uint16(spheroHead))
				time.Sleep(1000 * time.Millisecond)
				//spheroBot.SetRawMotorValues(sphero.Forward, cmdSpeed, sphero.Forward, cmdSpeed)
				spheroBot.SetRGB(255, 0, 0)
				//time.Sleep(1000 * time.Millisecond)
				spheroBot.Roll(0, uint16(spheroHead))
				sp <- Command{"stop", 0, 0, 0}
			case "go":
				spheroBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					spheroBot.Roll(DefaultRollSpeed, uint16(spheroHead))
					checkDuration(&it, "sphero")
				})
			case "spin":
				spheroBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					spheroBot.SetRawMotorValues(sphero.Forward, cmdSpeed, sphero.Reverse, cmdSpeed)
					if it%5 == 0 {
						spheroBot.SetRGB(uint8(gobot.Rand(255)),
							uint8(gobot.Rand(255)),
							uint8(gobot.Rand(255)))
					}
					checkDuration(&it, "sphero")
				})
			case "blink":
				ticker = gobot.Every(cmdInterval, func() {
					spheroBot.SetRGB(uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)))
					checkDuration(&it, "sphero")
				})
			case "boost":
				spheroBot.Boost(true)
			case "turn":
				spheroBot.SetRGB(0, 0, 255)
				spheroHead = (360 + (spheroHead + cmdDirection)) % 360
				spheroBot.Roll(0, uint16(spheroHead))
				time.Sleep(1 * time.Second)
				sp <- Command{"stop", 0, 0, 0}
			case "stop":
				if ticker != nil {
					ticker.Stop()
				}
				spheroBot.SetRGB(255, 0, 0)
				spheroBot.Roll(0, uint16(spheroHead))
				Complete <- true
			default:
				log.WithFields(log.Fields{"device": "sphero",
					"cmd": cmd,
				}).Error("Invalid command")
			}
		}
	}
	robot := gobot.NewRobot("spheroBot-"+port,
		[]gobot.Connection{adaptor},
		[]gobot.Device{spheroBot},
		work,
	)
	return robot
}
