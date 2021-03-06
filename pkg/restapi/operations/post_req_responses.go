// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/infracloudio/ollie-demo/pkg/models"
)

// PostReqOKCode is the HTTP code returned for type PostReqOK
const PostReqOKCode int = 200

/*PostReqOK success

swagger:response postReqOK
*/
type PostReqOK struct {

	/*
	  In: Body
	*/
	Payload *models.Resp `json:"body,omitempty"`
}

// NewPostReqOK creates PostReqOK with default headers values
func NewPostReqOK() *PostReqOK {
	return &PostReqOK{}
}

// WithPayload adds the payload to the post req o k response
func (o *PostReqOK) WithPayload(payload *models.Resp) *PostReqOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the post req o k response
func (o *PostReqOK) SetPayload(payload *models.Resp) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PostReqOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PostReqForbiddenCode is the HTTP code returned for type PostReqForbidden
const PostReqForbiddenCode int = 403

/*PostReqForbidden Forbidden

swagger:response postReqForbidden
*/
type PostReqForbidden struct {
}

// NewPostReqForbidden creates PostReqForbidden with default headers values
func NewPostReqForbidden() *PostReqForbidden {
	return &PostReqForbidden{}
}

// WriteResponse to the client
func (o *PostReqForbidden) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(403)
}

// PostReqMethodNotAllowedCode is the HTTP code returned for type PostReqMethodNotAllowed
const PostReqMethodNotAllowedCode int = 405

/*PostReqMethodNotAllowed Invalid input

swagger:response postReqMethodNotAllowed
*/
type PostReqMethodNotAllowed struct {
}

// NewPostReqMethodNotAllowed creates PostReqMethodNotAllowed with default headers values
func NewPostReqMethodNotAllowed() *PostReqMethodNotAllowed {
	return &PostReqMethodNotAllowed{}
}

// WriteResponse to the client
func (o *PostReqMethodNotAllowed) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(405)
}
