package ollieController

import (
	"fmt"
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
		ch <- OllieCommand{"stop", cmdDirection, 0, 0}
	}
}

const (
	defaultDir      = uint16(0)
	defaultSpeed    = uint8(255)
	defaultInterval = 100
	defaultDur      = 20000
)

var (
	cmdDirection uint16
	cmdSpeed     uint8
	cmdDuration  uint16
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
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.SetRawMotorValues(ollie.Forward, cmdSpeed, ollie.Forward, cmdSpeed)
					ollieBot.SetRGB(uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)))
					checkDuration(&it)
				})
			case "roll":
				ollieBot.Roll(0, cmdDirection)
				time.Sleep(1 * time.Second)
				ollieBot.SetRGB(0, 0, 255)
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.Roll(cmdSpeed, cmdDirection)
					checkDuration(&it)
				})
			case "spin":
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.SetRawMotorValues(ollie.Forward, cmdSpeed, ollie.Reverse, cmdSpeed)
					ollieBot.SetRGB(uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)),
						uint8(gobot.Rand(255)))
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
			case "stop":
				if ticker != nil {
					ticker.Stop()
				}
				ollieBot.SetRGB(255, 0, 0)
				ollieBot.Roll(0, cmdDirection)
			default:
				fmt.Println("invalid command")
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
