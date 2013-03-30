package image

const (
    ManifestVersion = 2
)

type ManifestState string

const (
    StateActive ManifestState = "active"
    StateUnactivated ManifestState = "unactivated"
    StateDisabled ManifestState = "disabled"
)

type Manifest struct {
    // Required
    V int `json:"v"`
    Uuid string `json:"uuid"`
    Owner string `json:"owner"`
    Name string `json:"name"`
    Version string `json:"version"`
    State ManifestState `json:"state"`
    Disabled bool `json:"disabled"`
    Public bool `json:"public"`
    PublishedAt string `json:"published_at"`
    Type string `json:"type"`
    Os string `json:"os"`
    Files []*ImageFile `json:"files"`

    // Required if type == zvol
    NicDriver string `json:"nic_driver"`
    DiskDriver string `json:"disk_driver"`
    CpuType string `json:"cpu_type"`
    ImageSize int64 `json:"image_size"`
    
    // Optional
    Description string `json:"description"`
}

type ImageFile struct {
    Sha1 string `json:"sha1"`
    Size int64 `json:"size"`
    Compression string `json:"compression"`
}
