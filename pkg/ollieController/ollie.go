package ollieController

import (
	"fmt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
	"gobot.io/x/gobot/platforms/sphero/ollie"
	"time"
)

func checkDuration(iter *int) {
	if *iter > 0 {
		*iter--
	}
	if *iter == 0 {
		ch <- OllieCommand{"stop", 0, 0, 0}
	}
}

// NewOllieBot Ollie bot
func NewOllieBot(port string) *gobot.Robot {
	cmdDirection := uint16(0)
	cmdSpeed := uint8(255)
	cmdDuration := 0
	cmdInterval := 200 * time.Millisecond
	bleAdaptor := ble.NewClientAdaptor(port)
	ollieBot := ollie.NewDriver(bleAdaptor)
	work := func() {
		ollieBot.AntiDOSOff()
		ollieBot.EnableStopOnDisconnect()
		var ticker *time.Ticker
		for {
			cmd := <-ch
			if cmd.Direction != 0 {
				cmdDirection = cmd.Direction
			}
			if cmd.Speed != 0 {
				cmdSpeed = cmd.Speed
			}
			if cmd.Duration > 0 {
				cmdDuration = int(cmd.Duration)
			} else {
				cmdDuration = -1000
			}
			// No of command iterations (1 per 200ms)
			it := cmdDuration / 200
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
				ticker = gobot.Every(cmdInterval, func() {
					ollieBot.Roll(cmdSpeed, cmdDirection)
					ollieBot.SetRGB(0, 0, 255)
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
				ollieBot.Stop()
				if ticker != nil {
					ollieBot.SetRGB(255, 0, 0)
					ticker.Stop()
				}
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
