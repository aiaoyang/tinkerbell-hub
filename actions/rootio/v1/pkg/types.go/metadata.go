package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Wrapper is a top level struct to unmarshal the metadata.
type Wrapper struct {
	Metadata Metadata `json:"metadata"`
}

// Metadata is an auto generated struct taken from a metadata request.
type Metadata struct {
	Instance Instance `json:"instance"`
}

// Instance is a dervied struct taken from a metadata request.
type Instance struct {
	CryptedRootPassword    string `json:"crypted_root_password"`
	Hostname               string `json:"hostname"`
	OperatingSystemVersion struct {
		Distro     string `json:"distro"`
		OsCodename string `json:"os_codename"`
		OsSlug     string `json:"os_slug"`
		Version    string `json:"version"`
	} `json:"operating_system_version"`
	Storage struct {
		Disks       []Disk       `json:"disks"`
		Filesystems []Filesystem `json:"filesystems"`
	} `json:"storage"`
}

// Filesystem defines the organisation of a filesystem.
type Filesystem struct {
	Mount struct {
		Create struct {
			Options []string `json:"options"`
		} `json:"create"`
		Device string `json:"device"`
		Format string `json:"format"`
		Point  string `json:"point"`
	} `json:"mount"`
}

// Disk defines the configuration for a disk.
type Disk struct {
	Device     string       `json:"device"`
	Partitions []Partitions `json:"partitions"`
	WipeTable  bool         `json:"wipe_table"`
}

// Partitions details the architecture.
type Partitions struct {
	Label  string `json:"label"`
	Number int    `json:"number"`
	Size   uint64 `json:"size"`
}

// RetrieveData retrieves metadata from Hegel.
func RetrieveData() (*Metadata, error) {
	metadataURL := os.Getenv("MIRROR_HOST")
	if metadataURL == "" {
		return nil, errors.New("unable to discover the metadata server from environment variable [MIRROR_HOST]")
	}

	metadataPort := os.Getenv("MIRROR_HOST_PORT")
	if metadataPort == "" {
		metadataPort = "5601"
	}
	metadataClient := http.Client{
		Timeout: time.Second * 60, // Timeout after 60 seconds (seems massively long is this dial-up?)
	}
	metaServerURL := fmt.Sprintf("http://%s:%s/metadata", metadataURL, metadataPort)
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, metaServerURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", metaServerURL, err)
	}

	req.Header.Set("User-Agent", "bootkit")

	res, getErr := metadataClient.Do(req)
	if getErr != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, err
	}

	var w Wrapper

	jsonErr := json.Unmarshal(body, &w)
	if jsonErr != nil {
		return nil, jsonErr
	}
	mdata := w.Metadata

	return &mdata, nil
}
