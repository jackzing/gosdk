package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"sync"

	json "github.com/json-iterator/go"
)

// definition of errors
var (
	ErrManifestNotExist    = errors.New("manifest not existed")
	ErrRecordNotExist      = errors.New("archive record not existed")
	ErrRecordStateNotExist = errors.New("archive record state not existed")
)

// Manifest represents all basic information of a snapshot.
type Manifest struct {
	Height         uint64 `json:"height"`
	Genesis        uint64 `json:"genesis"`
	BlockHash      string `json:"hash"`
	FilterID       string `json:"filterId"`
	MerkleRoot     string `json:"merkleRoot"`
	Namespace      string `json:"Namespace"`
	TxCount        uint64 `json:"txCount"`
	InvalidTxCount uint64 `json:"invalidTxCount,omitEmpty"`
	Status         uint   `json:"status"`
	DBVersion      string `json:"dbVersion"`
}

// Manifests defines a list of Manifest
type Manifests []Manifest

func (m Manifests) Len() int           { return len(m) }
func (m Manifests) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m Manifests) Less(i, j int) bool { return m[i].Height < m[j].Height }

/*
	Manifest manipulator
*/

// ManifestRW provides access to the read-write sets of manifest
type ManifestRW interface {
	Read(string) (Manifest, error)
	Write(Manifest) error
	List() (Manifests, error)
	Delete(string) error
	Contain(string) bool
	Search(uint64) (Manifest, error)
	Update(filterID string, mf Manifest) error
	SetStatus(filterID string, status uint) error
}

// ManifestHandler implements the ManifestRW interface
type ManifestHandler struct {
	filePath string
	lock     sync.RWMutex
	fsync    bool
}

// NewManifestHandler returns the handler of manifest according to given file path
func NewManifestHandler(fName string, fsync bool) *ManifestHandler {
	return &ManifestHandler{
		filePath: fName,
		fsync:    fsync,
	}
}

// Read reads the manifest according to given id
func (rwc *ManifestHandler) Read(id string) (Manifest, error) {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	buf, err := ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return Manifest{}, err
	}
	if len(buf) == 0 {
		return Manifest{}, ErrManifestNotExist
	}
	var manifests Manifests
	if err := json.Unmarshal(buf, &manifests); err != nil {
		return Manifest{}, err
	}
	for _, manifest := range manifests {
		if id == manifest.FilterID {
			return manifest, nil
		}
	}
	return Manifest{}, ErrManifestNotExist
}

// Write writes the given manifest
func (rwc *ManifestHandler) Write(manifest Manifest) error {
	rwc.lock.Lock()
	defer rwc.lock.Unlock()
	buf, _ := ioutil.ReadFile(rwc.filePath)
	var manifests Manifests
	if len(buf) != 0 {
		if err := json.Unmarshal(buf, &manifests); err != nil {
			return err
		}
	}
	manifests = append(manifests, manifest)
	if buf, err := json.MarshalIndent(manifests, "", "   "); err != nil {
		return err
	} else if err = rwc.writeBuf(buf); err != nil {
		return err
	}
	return nil
}

// List returns the list of manifests
func (rwc *ManifestHandler) List() (Manifests, error) {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	buf, err := ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return nil, err
	}
	if len(buf) == 0 {
		return nil, ErrManifestNotExist
	}
	var manifests Manifests
	if err := json.Unmarshal(buf, &manifests); err != nil {
		return nil, err
	}
	return manifests, nil
}

// Delete deletes the manifest according to given id
func (rwc *ManifestHandler) Delete(id string) error {
	rwc.lock.Lock()
	defer rwc.lock.Unlock()
	var deleted bool
	buf, err := ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return ErrManifestNotExist
	}
	var manifests Manifests
	if err := json.Unmarshal(buf, &manifests); err != nil {
		return err
	}
	for idx, manifest := range manifests {
		if manifest.FilterID == id {
			manifests = append(manifests[:idx], manifests[idx+1:]...)
			deleted = true
		}
	}
	if deleted {
		if buf, err := json.MarshalIndent(manifests, "", "   "); err != nil {
			return err
		} else if err := rwc.writeBuf(buf); err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return ErrManifestNotExist
	}
}

// Contain check if the manifest according to given id is exist or not
func (rwc *ManifestHandler) Contain(id string) bool {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	buf, err := ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return false
	}
	var manifests Manifests
	if err := json.Unmarshal(buf, &manifests); err != nil {
		return false
	}
	for _, manifest := range manifests {
		if manifest.FilterID == id {
			return true
		}
	}
	return false
}

// Search returns the manifest if its height is equal to the given height
func (rwc *ManifestHandler) Search(height uint64) (Manifest, error) {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	buf, err := ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return Manifest{}, err
	}
	if len(buf) == 0 {
		return Manifest{}, ErrManifestNotExist
	}
	var manifests Manifests
	err = json.Unmarshal(buf, &manifests)
	if err != nil {
		return Manifest{}, err
	}
	for _, manifest := range manifests {
		if manifest.Height == height {
			return manifest, err
		}
	}
	return Manifest{}, ErrManifestNotExist
}

// Update update the manifest field with specific filterID
func (rwc *ManifestHandler) Update(filterID string, mf Manifest) error {
	var (
		buf []byte
		err error
	)
	rwc.lock.Lock()
	defer rwc.lock.Unlock()
	buf, err = ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return ErrManifestNotExist
	}
	var manifests Manifests
	if err = json.Unmarshal(buf, &manifests); err != nil {
		return err
	}

	var i int
	for i = 0; i < manifests.Len(); i++ {
		if manifests[i].FilterID == filterID {
			break
		}
	}
	if i == manifests.Len() {
		return ErrManifestNotExist
	}
	manifests[i] = mf
	if buf, err = json.MarshalIndent(manifests, "", "   "); err != nil {
		return err
	}

	return rwc.writeBuf(buf)
}

// SetStatus set the manifest `status` field with true
func (rwc *ManifestHandler) SetStatus(filterID string, status uint) error {
	var (
		buf []byte
		err error
	)
	rwc.lock.Lock()
	defer rwc.lock.Unlock()
	buf, err = ioutil.ReadFile(rwc.filePath)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return ErrManifestNotExist
	}
	var manifests Manifests
	if err = json.Unmarshal(buf, &manifests); err != nil {
		return err
	}

	var i int
	for i = 0; i < manifests.Len(); i++ {
		if manifests[i].FilterID == filterID {
			break
		}
	}
	if i == manifests.Len() {
		return ErrManifestNotExist
	}
	manifests[i].Status = status
	if buf, err = json.MarshalIndent(manifests, "", "   "); err != nil {
		return err
	}

	return rwc.writeBuf(buf)
}

// CheckMetaExist check whether the snapshot.meta file exist
func (rwc *ManifestHandler) CheckMetaExist() bool {
	rwc.lock.RLock()
	defer rwc.lock.RUnlock()
	_, serr := os.Stat(rwc.filePath)
	if serr != nil && os.IsNotExist(serr) {
		return false
	}
	return true
}

func (rwc *ManifestHandler) writeBuf(buf []byte) error {
	file, err := os.OpenFile(rwc.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()
	_, err = file.WriteAt(buf, 0)
	if err != nil {
		return err
	}
	if rwc.fsync {
		return file.Sync()
	}
	return nil
}
