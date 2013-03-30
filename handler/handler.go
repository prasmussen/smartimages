package handler

import (
    "net/http"
    "github.com/prasmussen/smartimages/image"
    "github.com/prasmussen/smartimages/log"
    "github.com/prasmussen/smartimages/responder"
)

type Handler struct {
    images *image.Pool
    logger *log.Logger
}

func New(pool *image.Pool, logger *log.Logger) *Handler {
    return &Handler{
        images: pool,
        logger: logger,
    }
}

func (self *Handler) LogResponder(req *http.Request, res http.ResponseWriter) *LogResponder {
    return &LogResponder{
        Logger: self.logger.RequestStart(req),
        Responder: responder.New(res),
    }
}


func (self *Handler) GetImage() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.getImage(res, req, logres)
    }
}

func (self *Handler) GetImageFile() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.getImageFile(res, req, logres)
    }
}

func (self *Handler) ListImages() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.listImages(res, req, logres)
    }
}

func (self *Handler) CreateImage() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.createImage(res, req, logres)
    }
}

func (self *Handler) AddImageFile() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.addImageFile(res, req, logres)
    }
}

func (self *Handler) ImageAction() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.imageAction(res, req, logres)
    }
}

func (self *Handler) DeleteImage() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.deleteImage(res, req, logres)
    }
}

func (self *Handler) Ping() func(res http.ResponseWriter, req *http.Request) {
    return func(res http.ResponseWriter, req *http.Request) {
        logres := self.LogResponder(req, res)
        defer logres.Logger.RequestEnd()

        self.ping(res, req, logres)
    }
}
