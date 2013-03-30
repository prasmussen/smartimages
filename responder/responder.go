package responder

import (
    "strconv"
    "encoding/json"
    "encoding/base64"
    "net/http"
    "github.com/prasmussen/smartimages/errors"
)

type Responder struct {
    res http.ResponseWriter
}

func New(res http.ResponseWriter) *Responder {
    return &Responder{res}
}

func (self *Responder) JSON(v interface{}) {
    self.res.Header().Set("Content-Type", "application/json")
    json.NewEncoder(self.res).Encode(v)
}

func (self *Responder) Success(statusCode int) {
    self.res.WriteHeader(statusCode)
}

func (self *Responder) SetContentLength(length int64) {
    str := strconv.FormatInt(length, 10)
    self.res.Header().Set("Content-length", str)
}

func (self *Responder) SetContentMd5(md5sum []byte) {
    // Md5sum should be base64 encoded according to rfc1864
    b64 := base64.StdEncoding.EncodeToString(md5sum)
    self.res.Header().Set("Content-Md5", b64)
}

func (self *Responder) Error(err errors.Error) {
    self.res.WriteHeader(err.StatusCode())
    self.JSON(err.Data())
}
