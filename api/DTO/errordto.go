package dto

import (
	"encoding/json"
	"log"
	"time"
)

type ErrorResponse struct {
	Error string    `json:"error"`
	Time  time.Time `json:"time"`
}

func (e ErrorResponse) ToString() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(b)
}

func NewErrorDto(err error) ErrorResponse {
	res := ErrorResponse{
		Error: err.Error(),
		Time:  time.Now(),
	}
	return res
}
