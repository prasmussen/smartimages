package handler

import (
    "net/http"
)

var defaultPong = &Pong{
    Ping: "pong",
    Version: "1.0.0",
    Imgapi: true,
}

type Pong struct {
    Ping string `json:"ping"`
    Version string `json:"version"`
    Imgapi bool `json:"imgapi"`
}

func (self *Handler) ping(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    logres.JSON(defaultPong)
}
