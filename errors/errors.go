package errors

import (
    "fmt"
)

func ValidationFailed(err error) Error {
    return &e{"ValidationFailed", "Validation of parameters failed.", 422, err}
}

func InvalidParameter(err error) Error {
     return &e{"InvalidParameter", "Given parameter was invalid.", 422, err}
}

func ImageFilesImmutable(err error) Error {
    return &e{"ImageFilesImmutable", "Cannot modify files on an activated image.", 422, err}
}

func ImageAlreadyActivated(err error) Error {
    return &e{"ImageAlreadyActivated", "Image is already activated.", 422, err}
}

func NoActivationNoFile(err error) Error {
    return &e{"NoActivationNoFile", "Image must have a file to be activated.", 422, err}
}

func OperatorOnly(err error) Error {
    return &e{"OperatorOnly", "Operator-only endpoint called by a non-operator.", 403, err}
}

func ImageUuidAlreadyExists(err error) Error {
    return &e{"ImageUuidAlreadyExists", "Attempt to import an image with a conflicting UUID", 409, err}
}

func Upload(err error) Error {
    return &e{"Upload", "There was a problem with the upload.", 400, err}
}

func StorageIsDown(err error) Error {
    return &e{"StorageIsDown", "Storage system is down.", 503, err}
}

func InternalError(err error) Error {
    return &e{"InternalError", "Internal Server Error", 500, err}
}

func ResourceNotFound(err error) Error {
    return &e{"ResourceNotFound", "Not Found", 404, err}
}

func InvalidHeader(err error) Error {
    return &e{"InvalidHeader", "An invalid header was given in the request.", 400, err}
}

func ServiceUnavailableError(err error) Error {
    return &e{"ServiceUnavailableError", "Service Unavailable", 503, err}
}

func UnauthorizedError(err error) Error {
    return &e{"UnauthorizedError", "Unauthorized", 401, err}
}

func BadRequestError(err error) Error {
    return &e{"BadRequestError", "Bad Request", 400, err}
}

type Error interface {
    Data() *Data
    StatusCode() int
    Err() error
    Error() string
}

type Data struct {
    Code string `json:"code"`
    Description string `json:"description"`
}

type e struct {
    code string
    description string
    statusCode int
    err error
}

func (self *e) StatusCode() int {
    return self.statusCode
}

func (self *e) Data() *Data {
    return &Data{
        Code: self.code,
        Description: self.description,
    }
}

func (self *e) Err() error {
    return self.err
}

func (self *e) Error() string {
    if self.err == nil {
        return self.code
    }
    return fmt.Sprintf("External: %s, Internal: %s", self.code, self.err.Error())
}
