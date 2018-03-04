package handlers

import (
	"fmt"
	"io/ioutil"
	"strconv"

	middleware "github.com/go-openapi/runtime/middleware"
	bot "github.com/infracloudio/ollie-demo/pkg/botController"
	models "github.com/infracloudio/ollie-demo/pkg/models"
	operations "github.com/infracloudio/ollie-demo/pkg/restapi/operations"
	"github.com/masatana/go-textdistance"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Bot struct {
	Bot BotConf `yaml:"bot"`
}

type BotConf struct {
	Name  string   `yaml:"name"`
	IDs   []string `yaml:"ids"`
	Speed uint8    `yaml:"speed"`
}

var Config Bot
var validCmds [6]string

func dumpIntent(i *models.Intent) string {
	return fmt.Sprintf("%v: {command:%v direction:%v speed:%v duration:%v}", i.Name, i.Slots.Trick.Value, i.Slots.Direction.Value, i.Slots.Speed.Value, i.Slots.Duration.Value)
}

// ReadConfig from application.json
func ReadConfig() {
	configFile, err := ioutil.ReadFile("config/application.yaml")
	if err != nil {
		log.Fatal("application.yaml file read err", err)
	}
	err = yaml.Unmarshal(configFile, &Config)
	if err != nil {
		log.Info("error:", err)
	}
	bot.DefaultRollSpeed = Config.Bot.Speed
	validCmds = [6]string{"spin", "stop", "jump", "blink", "go", "turn"}
}

func buildResponse(title string, output string, repromptText string, shouldEndSession bool, command string, dir int16, speed uint8, dur uint16) models.Resp {
	outputSpeech := models.OutputSpeech{Type: "PlainText", Text: output}
	sessionAttr := models.Attributes{Command: command, Direction: dir, Speed: speed, Duration: dur}
	card := models.Card{Type: "Simple", Title: "SessionSpeechlet - " + title, Content: "SessionSpeechlet - " + output}
	reprompt := models.Reprompt{OutputSpeech: &models.OutputSpeech{Type: "PlainText", Text: repromptText}}
	return models.Resp{Response: &models.Response{OutputSpeech: &outputSpeech, Card: &card, Reprompt: &reprompt, ShouldEndSession: &shouldEndSession}, SessionAttributes: &sessionAttr, Version: "1.0"}
}

func getRetryResponse(cmd string) middleware.Responder {
	resp := buildResponse("Welcome", "Ollie doesn't understand "+cmd+". Please try again.", "What's next?", false, "", 0, 0, 0)
	r := operations.NewPostReqOK()
	r.Payload = &resp
	return r
}

func getWelcomeResponse() middleware.Responder {
	resp := buildResponse("Welcome", "Hello there! Please ask me to give Ollie commands like: spin or go", "What's next?", true, "", 0, 0, 0)
	r := operations.NewPostReqOK()
	r.Payload = &resp
	return r
}

func getLaunchResponse() middleware.Responder {
	resp := buildResponse("Welcome", "Connected", "What's next?", false, "", 0, 0, 0)
	r := operations.NewPostReqOK()
	r.Payload = &resp
	return r
}

func parserDirection(dir string) int16 {
	if dir == "" {
		return 0
	}
	// TODO: Add more directions
	switch dir {
	case "left":
		return -90
	case "straight":
		return 0
	case "right":
		return 90
	case "return":
		return 180
	case "reverse":
		return 180
	case "back":
		return 180
	default:
		return 0
	}
}

func getIntentResponse(req *models.Request, session *models.Session) middleware.Responder {
	var dir int16
	var speed uint8
	var dur uint16
	found := false
	endSession := false
	defRespText := "what's next?"
	cmd := req.Intent.Slots.Trick.Value
	log.WithFields(log.Fields{
		"command": cmd,
	}).Info("Command received by backend")
	if cmd == "" {
		log.Info("Command is empty. Sending Welcome response")
		return getWelcomeResponse()
	} else {
		for _, c := range validCmds {
			// check JaroWinklerDistance
			dist := textdistance.JaroWinklerDistance(c, cmd)
			if dist == 1 {
				found = true
				break
			} else if dist >= 0.70 {
				log.Info("Command predicted: " + cmd)
				found = true
				cmd = c
				break
			}
		}
	}

	if found == false {
		log.Info("Invalid command " + cmd + ". Sending retry response")
		return getRetryResponse(cmd)
	}

	if cmd == "stop" {
		endSession = true
		defRespText = ""
	}
	dir = parserDirection(req.Intent.Slots.Direction.Value)
	if s, err := strconv.ParseUint(req.Intent.Slots.Speed.Value, 10, 8); err == nil {
		speed = uint8(s)
	} else {
		speed = 0
	}
	if req.Intent.Slots.Duration.Value == "half" {
		dur = 500
	} else if s, err := strconv.ParseUint(req.Intent.Slots.Duration.Value, 10, 16); err == nil {
		dur = uint16(s) * 1000
	} else {
		dur = 0
	}
	resp := buildResponse("Welcome", defRespText, "What's next?", endSession, cmd, dir, speed, dur)
	sendCommandToBot(resp)
	r := operations.NewPostReqOK()
	r.Payload = &resp
	return r
}

func onLaunch(req *models.Request, session *models.Session) middleware.Responder {
	return getLaunchResponse()
}

func onIntent(req *models.Request, session *models.Session) middleware.Responder {
	if req.Intent.Name == "Command" {
		return getIntentResponse(req, session)
	} else {
		return getWelcomeResponse()
	}
}

// Send command to ollieController channel
func sendCommandToBot(resp models.Resp) {
	c := resp.SessionAttributes.Command
	dir := resp.SessionAttributes.Direction
	speed := resp.SessionAttributes.Speed
	// Convert dur to milliseconds
	dur := resp.SessionAttributes.Duration
	cmd := bot.Command{Command: c, Direction: dir, Speed: speed, Duration: dur}
	log.WithFields(log.Fields{
		"command": cmd,
	}).Info("Sending command to " + Config.Bot.Name)
	if Config.Bot.Name == "ollie" {
		bot.SendCommandToOllie(cmd)
	} else if Config.Bot.Name == "sphero" {
		bot.SendCommandToSphero(cmd)
	}

	// wait for ack
	_ = <-bot.Complete
}

// Handle POST requests
func HandlePostReq(req *models.Req) middleware.Responder {
	if req.Request.Type == "LaunchRequest" {
		log.WithFields(log.Fields{
			"intent": req.Request.Type,
		}).Info("Received request")
		return onLaunch(req.Request, req.Session)
	} else if req.Request.Type == "IntentRequest" {
		log.WithFields(log.Fields{
			"intentType": req.Request.Type,
			"intent":     dumpIntent(req.Request.Intent),
		}).Info("Received request")
		return onIntent(req.Request, req.Session)
	} else {
		return onLaunch(req.Request, req.Session)
	}

}
