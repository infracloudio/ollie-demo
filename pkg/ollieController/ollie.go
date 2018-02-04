package ollieController

import (
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/ollie"
	"time"
)

func checkDuration(iter *uint16) {
	if *iter > 0 {
		*iter--
	}
	if *iter == 0 {
		ch <- OllieCommand{"stop", 0, 0, 0}
	}
}

const (
	defaultDir       = 0
	defaultSpeed     = uint8(255)
	defaultRollSpeed = uint8(100)
	defaultInterval  = 100
	defaultDur       = 5000
)

var (
	cmdDirection int16
	cmdSpeed     uint8
	cmdDuration  uint16
	ollieHead    int16
)

// NewOllieBot Ollie bot
func NewOllieBot(port string) *gobot.Robot {
	cmdDirection = defaultDir
	cmdSpeed = defaultSpeed
	cmdDuration = defaultDur
	cmdInterval := defaultInterval * time.Millisecond
	bleAdaptor := ble.NewClientAdaptor(port)
	ollieBot := ollie.NewDriver(bleAdaptor)
	work := func() {
		ollieHead = 0
		ollieBot.AntiDOSOff()
		ollieBot.EnableStopOnDisconnect()
		ollieBot.SetStabilization(false)
		ollieBot.Roll(0, 0)
		var ticker *time.Ticker
		for {
			cmd := <-ch
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
				cmdDuration = defaultDur
			}
			// No of command iterations (1 per 200ms)
			it := cmdDuration / defaultInterval
			switch cmd.Command {
			case "jump":
				ollieBot.SetRGB(0, 0, 255)
				ollieBot.SetRawMotorValues(ollie.Forward, cmdSpeed, ollie.Forward, cmdSpeed)
				ollieBot.SetRGB(255, 0, 0)
				ollieBot.Roll(0, uint16(ollieHead))
			case "go":
				ollieBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.Roll(defaultRollSpeed, uint16(ollieHead))
					checkDuration(&it)
				})
			case "spin":
				ollieBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.SetRawMotorValues(ollie.Forward, cmdSpeed, ollie.Reverse, cmdSpeed)
					if it % 5 == 0 {
						ollieBot.SetRGB(uint8(gobot.Rand(255)),
							uint8(gobot.Rand(255)),
							uint8(gobot.Rand(255)))
					}
					checkDuration(&it)
				})
			case "blink":
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.SetRGB(uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)))
					checkDuration(&it)
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
				log.WithFields(log.Fields{
					"cmd": cmd,
				}).Info("Invalid command received")
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
