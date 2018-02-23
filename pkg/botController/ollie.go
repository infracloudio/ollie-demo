package botController

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/ollie"
	"time"

	log "github.com/sirupsen/logrus"
)

// NewOllieBot Ollie bot
func NewOllieBot(port string) *gobot.Robot {
	cmdDirection = defaultDir
	cmdSpeed = defaultSpeed
	cmdDuration = DefaultDur
	cmdInterval := defaultInterval * time.Millisecond
	bleAdaptor := ble.NewClientAdaptor(port)
	ollieBot := ollie.NewDriver(bleAdaptor)
	work := func() {
		ollieHead := int16(0)
		ollieBot.AntiDOSOff()
		ollieBot.EnableStopOnDisconnect()
		ollieBot.SetStabilization(false)
		ollieBot.Roll(0, 0)
		var ticker *time.Ticker
		for {
			cmd := <-ol
			log.WithFields(log.Fields{
				"device": "ollie",
				"cmd":    cmd,
			}).Info("command received by ollie driver")
			ollieBot.Wake()
			// Stop previous command execution
			if ticker != nil {
				ollieBot.SetRGB(255, 0, 0)
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
				ollieBot.SetRGB(0, 0, 255)
				ollieBot.Roll(255, uint16(ollieHead))
				time.Sleep(1000 * time.Millisecond)
				ollieBot.SetRawMotorValues(ollie.Forward, cmdSpeed, ollie.Forward, cmdSpeed)
				ollieBot.SetRGB(255, 0, 0)
				time.Sleep(1000 * time.Millisecond)
				ollieBot.Roll(0, uint16(ollieHead))
			case "go":
				ollieBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.Roll(DefaultRollSpeed, uint16(ollieHead))
					checkDuration(&it, "ollie")
				})
			case "spin":
				ollieBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.SetRawMotorValues(ollie.Forward, cmdSpeed, ollie.Reverse, cmdSpeed)
					if it%5 == 0 {
						ollieBot.SetRGB(uint8(gobot.Rand(255)),
							uint8(gobot.Rand(255)),
							uint8(gobot.Rand(255)))
					}
					checkDuration(&it, "ollie")
				})
			case "blink":
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.SetRGB(uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)))
					checkDuration(&it, "ollie")
				})
			case "boost":
				ollieBot.Boost(true)
			case "turn":
				ollieHead = (360 + (ollieHead + cmdDirection)) % 360
				ollieBot.Roll(0, uint16(ollieHead))
			case "stop":
				if ticker != nil {
					ticker.Stop()
				}
				ollieBot.SetRGB(255, 0, 0)
				ollieBot.Roll(0, uint16(ollieHead))
			default:
				log.WithFields(log.Fields{"device": "ollie",
					"cmd": cmd,
				}).Error("Invalid command")
			}
		}
	}
	robot := gobot.NewRobot("ollieBot-"+port,
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ollieBot},
		work,
	)
	return robot
}
