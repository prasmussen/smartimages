package handler

import (
    "github.com/prasmussen/smartimages/log"
    "github.com/prasmussen/smartimages/responder"
    "github.com/prasmussen/smartimages/errors"
)

type LogResponder struct {
    Logger *log.RequestLogger
    Responder *responder.Responder
}

func (self *LogResponder) JSON(i interface{}) {
    self.Responder.JSON(i)
    self.Logger.Success()
}

func (self *LogResponder) Success(statusCode int) {
    self.Responder.Success(statusCode)
    self.Logger.Success()
}

func (self *LogResponder) Error(err errors.Error) {
    self.Responder.Error(err)    
    self.Logger.Error(err)    
}
