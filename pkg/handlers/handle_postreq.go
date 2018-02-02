package handlers

import (
	"fmt"
	middleware "github.com/go-openapi/runtime/middleware"
	models "github.com/infracloudio/ollie-demo/pkg/models"
	ollieBot "github.com/infracloudio/ollie-demo/pkg/ollieController"
	operations "github.com/infracloudio/ollie-demo/pkg/restapi/operations"
	"strconv"
)

func buildResponse(title string, output string, repromptText string, shouldEndSession bool, command string, dir int16, speed uint8, dur uint16) models.Resp {
	outputSpeech := models.OutputSpeech{Type: "PlainText", Text: output}
	sessionAttr := models.Attributes{Command: command, Direction: dir, Speed: speed, Duration: dur}
	card := models.Card{Type: "Simple", Title: "SessionSpeechlet - " + title, Content: "SessionSpeechlet - " + output}
	reprompt := models.Reprompt{OutputSpeech: &models.OutputSpeech{Type: "PlainText", Text: repromptText}}
	return models.Resp{Response: &models.Response{OutputSpeech: &outputSpeech, Card: &card, Reprompt: &reprompt, ShouldEndSession: &shouldEndSession}, SessionAttributes: &sessionAttr, Version: "1.0"}
}

func getWelcomeResponse() middleware.Responder {
	resp := buildResponse("Welcome", "Hello there! Please ask me to give Ollie commands like: Tell ollie to spin", "What's next?", true, "", 0, 0, 0)
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

func getIntentResponse(req *models.Request) middleware.Responder {
	var dir int16
	var speed uint8
	var dur uint16
	cmd := req.Intent.Slots.Trick.Value
	if cmd == "" {
		return getWelcomeResponse()
	}
	dir = parserDirection(req.Intent.Slots.Direction.Value)
	if s, err := strconv.ParseUint(req.Intent.Slots.Speed.Value, 10, 8); err == nil {
		speed = uint8(s)
	} else {
		speed = 0
	}
	if s, err := strconv.ParseUint(req.Intent.Slots.Duration.Value, 10, 16); err == nil {
		dur = uint16(s)
	} else {
		dur = 0
	}
	resp := buildResponse("Welcome", "", "What's next?", true, cmd, dir, speed, dur)
	sendCommandToOllie(resp)
	r := operations.NewPostReqOK()
	r.Payload = &resp
	return r
}

func onLaunch(req *models.Request, session *models.Session) middleware.Responder {
	return getWelcomeResponse()
}

func onIntent(req *models.Request, session *models.Session) middleware.Responder {
	if req.Intent.Name == "OllieCommand" {
		return getIntentResponse(req)
	} else {
		return getWelcomeResponse()
	}
}

// Send command to ollieController channel
func sendCommandToOllie(resp models.Resp) {
	// Wait till Alexa gets the response
	c := resp.SessionAttributes.Command
	dir := resp.SessionAttributes.Direction
	speed := resp.SessionAttributes.Speed
	// Convert dur to milliseconds
	dur := resp.SessionAttributes.Duration * 1000
	cmd := ollieBot.OllieCommand{Command: c, Direction: dir, Speed: speed, Duration: dur}
	fmt.Println("Command - ", cmd)
	ollieBot.SendCommand(cmd)
}

// Handle POST requests
func HandlePostReq(req *models.Req) middleware.Responder {
	fmt.Println(req.Request.Type)
	if req.Request.Type == "LaunchRequest" {
		return onLaunch(req.Request, req.Session)
	} else if req.Request.Type == "IntentRequest" {
		return onIntent(req.Request, req.Session)
	} else {
		return onLaunch(req.Request, req.Session)
	}

}
