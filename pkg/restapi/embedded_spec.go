// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Alexa custom skill service endpoints for Ollie",
    "title": "ollie-skill-api",
    "version": "1.0.0"
  },
  "host": "localhost:5000",
  "paths": {
    "/": {
      "post": {
        "description": "Alexa service request",
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "operationId": "postReq",
        "parameters": [
          {
            "description": "Request body",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Req"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "success",
            "schema": {
              "$ref": "#/definitions/Resp"
            }
          },
          "403": {
            "description": "Forbidden"
          },
          "405": {
            "description": "Invalid input"
          }
        }
      }
    }
  },
  "definitions": {
    "Application": {
      "type": "object",
      "properties": {
        "applicationId": {
          "type": "string"
        }
      }
    },
    "Attributes": {
      "type": "object",
      "properties": {
        "command": {
          "type": "string"
        },
        "direction": {
          "type": "integer",
          "format": "uint16"
        },
        "duration": {
          "type": "integer",
          "format": "uint16"
        },
        "speed": {
          "type": "integer",
          "format": "uint8"
        }
      }
    },
    "AudioPlayer": {
      "type": "object",
      "properties": {
        "offsetInMilliseconds": {
          "type": "integer",
          "format": "int32"
        },
        "playerActivity": {
          "description": "Indicates the last known state of audio",
          "type": "string",
          "enum": [
            "IDLE",
            "PAUSED",
            "PLAYING",
            "BUFFER_UNDERRUN",
            "FINISHED",
            "STOPPED"
          ]
        },
        "token": {
          "type": "string"
        }
      }
    },
    "Card": {
      "type": "object",
      "properties": {
        "content": {
          "type": "string"
        },
        "image": {
          "$ref": "#/definitions/Image"
        },
        "text": {
          "type": "string"
        },
        "title": {
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      }
    },
    "Context": {
      "type": "object",
      "properties": {
        "AudioPlayer": {
          "$ref": "#/definitions/AudioPlayer"
        },
        "system": {
          "$ref": "#/definitions/System"
        }
      }
    },
    "Device": {
      "type": "object",
      "properties": {
        "deviceId": {
          "type": "string"
        },
        "supportedInterfaces": {
          "$ref": "#/definitions/SupportedIntf"
        }
      }
    },
    "Image": {
      "type": "object",
      "properties": {
        "largeImageUrl": {
          "type": "string"
        },
        "smallImageUrl": {
          "type": "string"
        }
      }
    },
    "Intent": {
      "type": "object",
      "properties": {
        "confirmationStatus": {
          "type": "string"
        },
        "name": {
          "type": "string"
        },
        "slots": {
          "$ref": "#/definitions/Slots"
        }
      }
    },
    "OutputSpeech": {
      "type": "object",
      "properties": {
        "text": {
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      }
    },
    "Permissions": {
      "type": "object",
      "properties": {
        "consentToken": {
          "type": "string"
        }
      }
    },
    "Reprompt": {
      "type": "object",
      "properties": {
        "outputSpeech": {
          "$ref": "#/definitions/OutputSpeech"
        }
      }
    },
    "Req": {
      "type": "object",
      "required": [
        "session",
        "context",
        "request"
      ],
      "properties": {
        "context": {
          "$ref": "#/definitions/Context"
        },
        "request": {
          "$ref": "#/definitions/Request"
        },
        "session": {
          "$ref": "#/definitions/Session"
        },
        "version": {
          "type": "string",
          "example": "1.0"
        }
      }
    },
    "Request": {
      "type": "object",
      "properties": {
        "dialogState": {
          "type": "string"
        },
        "intent": {
          "$ref": "#/definitions/Intent"
        },
        "locale": {
          "type": "string"
        },
        "requestId": {
          "type": "string"
        },
        "timestamp": {
          "type": "string"
        },
        "type": {
          "type": "string"
        }
      }
    },
    "Resp": {
      "type": "object",
      "properties": {
        "response": {
          "$ref": "#/definitions/Response"
        },
        "sessionAttributes": {
          "$ref": "#/definitions/Attributes"
        },
        "version": {
          "type": "string"
        }
      }
    },
    "Response": {
      "type": "object",
      "properties": {
        "card": {
          "$ref": "#/definitions/Card"
        },
        "outputSpeech": {
          "$ref": "#/definitions/OutputSpeech"
        },
        "reprompt": {
          "$ref": "#/definitions/Reprompt"
        },
        "shouldEndSession": {
          "type": "boolean",
          "default": true
        }
      }
    },
    "Session": {
      "type": "object",
      "properties": {
        "application": {
          "$ref": "#/definitions/Application"
        },
        "attributes": {
          "$ref": "#/definitions/Attributes"
        },
        "new": {
          "type": "boolean",
          "default": true
        },
        "sessionId": {
          "type": "string"
        },
        "user": {
          "$ref": "#/definitions/User"
        }
      }
    },
    "SlotName": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "value": {
          "type": "string"
        }
      }
    },
    "Slots": {
      "type": "object",
      "properties": {
        "direction": {
          "$ref": "#/definitions/SlotName"
        },
        "duration": {
          "$ref": "#/definitions/SlotName"
        },
        "speed": {
          "$ref": "#/definitions/SlotName"
        },
        "trick": {
          "$ref": "#/definitions/SlotName"
        }
      }
    },
    "SupportedIntf": {
      "type": "object",
      "properties": {
        "AudioPlayer": {
          "$ref": "#/definitions/AudioPlayer"
        }
      }
    },
    "System": {
      "type": "object",
      "properties": {
        "apiAccessToken": {
          "type": "string"
        },
        "apiEndpoint": {
          "type": "string"
        },
        "application": {
          "$ref": "#/definitions/Application"
        },
        "device": {
          "$ref": "#/definitions/Device"
        },
        "user": {
          "$ref": "#/definitions/User"
        }
      }
    },
    "User": {
      "type": "object",
      "properties": {
        "accessToken": {
          "type": "string"
        },
        "permissions": {
          "$ref": "#/definitions/Permissions"
        },
        "userId": {
          "type": "string"
        }
      }
    }
  }
}`))
}
