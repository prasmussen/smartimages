package handler

import (
    "net/http"
    "io"
    "encoding/json"
    "github.com/prasmussen/smartimages/image"
    "github.com/prasmussen/smartimages/errors"
)

func (self *Handler) getImage(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    query := req.URL.Query()

    uuid := query.Get(":uuid")

    manifest, err := self.images.Get(uuid)
    if err != nil {
        logres.Error(err)
        return
    }

    logres.JSON(manifest)
}

func (self *Handler) getImageFile(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    query := req.URL.Query()

    uuid := query.Get(":uuid")

    reader, metadata, err := self.images.GetFile(uuid)
    if err != nil {
        logres.Error(err)
        return
    }

    defer reader.Close()

    // Set content-length and content-md5 header required by imgadm.
    // The data is sent as chunked-encoding and content-length is not
    // really needed, but the imgadm expects it
    logres.Responder.SetContentMd5(metadata.Md5sum)
    logres.Responder.SetContentLength(metadata.Size)

    // Write file to response
    io.Copy(res, reader)

    logres.Logger.Success()
}

func (self *Handler) listImages(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    query := req.URL.Query()

    filters := make([]image.Filter, 0)

    for key, values := range query {
        filter, ok := image.GetFilter(key, values[0])
        if !ok {
            continue
        }
        filters = append(filters, filter)
    }

    // Filter on active images by default if no state was explicitly set
    if query.Get("state") == "" {
        filters = append(filters, image.StateFilter("active"))
    }

    manifests := self.images.List(filters)
    logres.JSON(manifests)
}

func (self *Handler) createImage(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    manifest := &image.Manifest{}
    if err := json.NewDecoder(req.Body).Decode(manifest); err != nil {
        logres.Error(errors.InternalError(err))
        return
    }

    if err := self.images.Create(manifest); err != nil {
        logres.Error(err)
        return
    }

    logres.JSON(manifest)
}

func (self *Handler) addImageFile(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    query := req.URL.Query()

    // Close body independent of the outcome
    defer req.Body.Close()

    // Grab uuid
    uuid := query.Get(":uuid")

    // Grab compression type
    compression := query.Get("compression")
    if compression == "" {
        logres.Error(errors.InvalidParameter(nil))
        return
    }

    manifest, err := self.images.AddFile(uuid, compression, req.Body)
    if err != nil {
        logres.Error(err)
        return
    }

    logres.JSON(manifest)
}

func (self *Handler) imageAction(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    query := req.URL.Query()

    // Close body
    defer req.Body.Close()

    uuid := query.Get(":uuid")
    action := query.Get("action")

    var manifest *image.Manifest
    var err errors.Error

    switch action {
    case "activate":
        manifest, err = self.images.Activate(uuid)        
    case "enable":
        manifest, err = self.images.SetDisabled(uuid, false)
    case "disable":
        manifest, err = self.images.SetDisabled(uuid, true)
    case "":
        err = errors.InvalidParameter(nil)
    default:
        err = errors.InvalidParameter(nil)
    }

    if err != nil {
        logres.Error(err)
        return
    }

    logres.JSON(manifest)
}

func (self *Handler) deleteImage(res http.ResponseWriter, req *http.Request, logres *LogResponder) {
    query := req.URL.Query()

    // Grab uuid
    uuid := query.Get(":uuid")

    if err := self.images.Delete(uuid); err != nil {
        logres.Error(err)
        return
    }

    logres.Success(204)
}
