package image

import (
    "fmt"
    "time"
    "io"
    "sync"
    "os"
    "encoding/json"
    "path/filepath"
    "crypto/sha1"
    "crypto/md5"
    "io/ioutil"
    "code.google.com/p/go-uuid/uuid"
    "github.com/prasmussen/smartimages/errors"
)

const (
    ManifestsFname = "manifests.json"
)

var FileExtensions = map[string]string{
    "bzip2": "bz2",
    "gzip": "gz",
    "none": "raw",
}

type FileMetadata struct {
    Md5sum []byte
    Size int64
}

type Pool struct {
    imageDir string
    manifests []*Manifest
    mutex *sync.Mutex
}

func NewImagePool(imageDir string) *Pool {
    return &Pool{
        imageDir: imageDir,
        manifests: loadManifests(),
        mutex: &sync.Mutex{},
    }
}

func loadManifests() []*Manifest {
    manifests := make([]*Manifest, 0)

    f, err := os.Open(ManifestsFname)
    if err != nil {
        return manifests
    }

    defer f.Close()

    json.NewDecoder(f).Decode(&manifests)
    return manifests
}

func saveManifests(manifests []*Manifest) error {
    // Grab a temp file
    f, err := ioutil.TempFile(".", "." + ManifestsFname)
    if err != nil {
        return err
    }

    // Close temp file on function exit its not already closed
    defer f.Close()

    // Write manifests to temp file
    if err := json.NewEncoder(f).Encode(manifests); err != nil {
        return err
    }

    // Close temp file
    f.Close()

    // Overwrite the old manifests file with the new one
    // Which is an atomic operation on sane OS's
    if err := os.Rename(f.Name(), ManifestsFname); err != nil {
        return err
    }

    return nil
}

func (self *Pool) Get(uuid string) (*Manifest, errors.Error) {
    manifest, ok := self.findManifest(uuid)
    if !ok {
        return nil, errors.ResourceNotFound(nil)
    }

    return manifest, nil
}

func (self *Pool) GetFile(uuid string) (io.ReadCloser, *FileMetadata, errors.Error) {
    manifest, ok := self.findManifest(uuid)
    if !ok {
        err := fmt.Errorf("Manifest not found")
        return nil, nil, errors.ResourceNotFound(err)
    }

    if len(manifest.Files) == 0 {
        err := fmt.Errorf("Manifest has no files")
        return nil, nil, errors.ResourceNotFound(err)
    }

    // Find image extension by compression type
    compression := manifest.Files[0].Compression
    ext := FileExtensions[compression]

    // Find absolute path for image and md5file
    imageFpath := filepath.Join(self.imageDir, fmt.Sprintf("%s.%s", uuid, ext))
    md5Fpath := filepath.Join(self.imageDir, uuid + ".md5")

    // Read md5 sum from file
    md5sum, err := ioutil.ReadFile(md5Fpath)
    if err != nil {
        return nil, nil, errors.InternalError(err)
    }

    // Open file
    f, err := os.Open(imageFpath)
    if err != nil {
        return nil, nil, errors.InternalError(err)
    }

    metadata := &FileMetadata{
        Md5sum: md5sum,
        Size: manifest.Files[0].Size,
    }

    return f, metadata, nil
}

func (self *Pool) List(filters []Filter) []*Manifest {
    self.lock()
    defer self.unlock()

    manifests := make([]*Manifest, 0)

    for _, m := range self.manifests {
        if MatchManifest(filters, m) {
            manifests = append(manifests, m)
        }
    }

    return manifests
}

func (self *Pool) Create(m *Manifest) errors.Error {
    m.V = ManifestVersion
    m.Uuid = uuid.NewUUID().String()
    m.State = StateUnactivated
    m.Disabled = true
    m.Public = true
    m.Files = make([]*ImageFile, 0)

    if err := self.addManifest(m); err != nil {
        return errors.InternalError(err)
    }

    return nil
}

func (self *Pool) Delete(uuid string) errors.Error {
    // Find manifest with matching uuid
    _, ok := self.findManifest(uuid)
    if !ok {
        return errors.ResourceNotFound(nil)
    }

    self.lock()

    // Remove manifest from internal slice first
    manifests := make([]*Manifest, 0)
    for _, m := range self.manifests {
        if m.Uuid != uuid {
            manifests = append(manifests, m)
        }
    }
    self.manifests = manifests

    // Save manifests to disk
    if err := saveManifests(self.manifests); err != nil {
        return errors.InternalError(err)
    }

    // No need to keep lock anymore
    self.unlock()

    // Find all files starting with the uuid of the image
    pattern := filepath.Join(self.imageDir, uuid + ".*")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return errors.InternalError(err)
    }

    // Delete files
    for _, fname := range matches {
        os.Remove(fname)
    }

    return nil
}

func (self *Pool) AddFile(uuid, compression string, reader io.Reader) (*Manifest, errors.Error) {
    // Find manifest with matching uuid
    manifest, ok := self.findManifest(uuid)
    if !ok {
        return nil, errors.ResourceNotFound(nil)
    }

    // Make sure the manifest has the correct state
    if manifest.State != StateUnactivated {
        return nil, errors.ImageAlreadyActivated(nil)
    }

    // Resolve file extension for the given compression type
    ext, ok := FileExtensions[compression]
    if !ok {
        return nil, errors.InvalidParameter(nil)
    }

    // Find absolute path for image and md5file
    imageFpath := filepath.Join(self.imageDir, fmt.Sprintf("%s.%s", uuid, ext))
    md5Fpath := filepath.Join(self.imageDir, uuid + ".md5")

    // Create destination directory if it does not exist
    err := os.MkdirAll(self.imageDir, 0775)
    if err != nil {
        return nil, errors.InternalError(err)
    }

    // Open image file
    f, err := os.Create(imageFpath)
    if err != nil {
        return nil, errors.InternalError(err)
    }

    // Remember to close files
    defer f.Close()

    // Calcluate sha1 and md5 sum while writing image to disk
    // Sha1 is a required field in the manifest
    shaHash := sha1.New()
    sha1Reader := io.TeeReader(reader, shaHash)

    // Md5 is needed by imgadm client which expects the
    // content-md5 header to be present
    md5Hash := md5.New()
    md5Reader := io.TeeReader(sha1Reader, md5Hash)
    
    // Write image to disk
    nBytes, err := io.Copy(f, md5Reader)
    if err != nil {
        return nil, errors.Upload(err)
    }

    // Write md5sum to file
    md5sum := md5Hash.Sum(nil)
    if err := ioutil.WriteFile(md5Fpath, md5sum, 0660); err != nil {
        return nil, errors.InternalError(err)
    }

    // Add image file to manifest
    imageFile := &ImageFile{
        Sha1: fmt.Sprintf("%x", shaHash.Sum(nil)),
        Compression: compression,
        Size: nBytes,
    }

    // Update manifest
    self.lock()
    defer self.unlock()
    manifest.Files = []*ImageFile{imageFile}

    // Save manifests to disk
    if err := saveManifests(self.manifests); err != nil {
        return nil, errors.InternalError(err)
    }

    return manifest, nil
}

func (self *Pool) Activate(uuid string) (*Manifest, errors.Error) {
    // Find manifest with the given uuid
    manifest, ok := self.findManifest(uuid)
    if !ok {
        return nil, errors.ResourceNotFound(nil)
    }

    // Make sure that an image file has been uploaded
    if len(manifest.Files) == 0 {
        return nil, errors.NoActivationNoFile(nil)
    }

    // Make sure it has not been activated before
    if manifest.State != StateUnactivated {
        return nil, errors.ImageAlreadyActivated(nil)
    }

    self.lock()
    defer self.unlock()

    // Activate image
    manifest.State = StateActive
    manifest.Disabled = false
    manifest.PublishedAt = time.Now().Format(time.RFC3339)

    // Save manifests to disk
    if err := saveManifests(self.manifests); err != nil {
        return nil, errors.InternalError(err)
    }

    return manifest, nil
}

func (self *Pool) SetDisabled(uuid string, disabled bool) (*Manifest, errors.Error) {
    // Find manifest with the given uuid
    manifest, ok := self.findManifest(uuid)
    if !ok {
        return nil, errors.ResourceNotFound(nil)
    }

    if !disabled && manifest.State == StateUnactivated {
        // Image must be activated before it can be enabled
        return nil, errors.ServiceUnavailableError(nil)
    }

    self.lock()
    defer self.unlock()

    // Enable / disable the image
    if disabled {
        manifest.State = StateDisabled
    } else {
        manifest.State = StateActive
    }

    manifest.Disabled = disabled

    // Save manifests to disk
    if err := saveManifests(self.manifests); err != nil {
        return nil, errors.InternalError(err)
    }

    return manifest, nil
}

func (self *Pool) findManifest(uuid string) (*Manifest, bool) {
    self.lock()
    defer self.unlock()

    for _, m := range self.manifests {
        if m.Uuid == uuid {
            return m, true
        }
    }

    return nil, false
}

func (self *Pool) addManifest(m *Manifest) error {
    self.lock()
    defer self.unlock()

    manifests := append(self.manifests, m)

    // Save manifests to disk
    if err := saveManifests(manifests); err != nil {
        return err
    }

    // Update internal manifests slice
    self.manifests = manifests

    return nil
}

func (self *Pool) lock() {
    self.mutex.Lock()
}

func (self *Pool) unlock() {
    self.mutex.Unlock()
}
