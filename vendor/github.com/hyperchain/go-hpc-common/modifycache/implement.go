package modifycache

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/hyperchain/go-hpc-common/math"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"sync"
)

// modifyCacheImpl is the implement of ModifyCache
type modifyCacheImpl struct {
	m sync.Map
}

// NewModifyCache generate a modify cache (mainly for executor)
func NewModifyCache() ModifyCache {
	return &modifyCacheImpl{}
}

// Set sets modify cache
func (m *modifyCacheImpl) Set(set ModifySet) {
	m.m.Store(set.GetSeqNo(), set)
}

// Get returns modify set with given sequence number
func (m *modifyCacheImpl) Get(seqNo uint64) (ModifySet, bool) {
	v, exist := m.m.Load(seqNo)
	if !exist {
		return nil, false
	}
	return v.(ModifySet), true
}

// Del removes modify set with given sequence number
func (m *modifyCacheImpl) Del(seqNo uint64) {
	m.m.Delete(seqNo)
}

// TrimBefore remove modify set before SeqNo (include SeqNo)
func (m *modifyCacheImpl) TrimBefore(seqNo uint64) {
	var toDelete []uint64
	m.m.Range(func(key, _ interface{}) bool {
		storedSeqNo := key.(uint64)
		if storedSeqNo <= seqNo {
			toDelete = append(toDelete, storedSeqNo)
		}
		return true
	})
	for _, d := range toDelete {
		m.m.Delete(d)
	}
}

// Reset clear all modify sets
func (m *modifyCacheImpl) Reset() {
	m.TrimBefore(math.MaxInt64)
}

// Len returns the length of modify cache
func (m *modifyCacheImpl) Len() int {
	var counter int
	m.m.Range(func(key, _ interface{}) bool {
		counter++
		return true
	})
	return counter
}

// ModifySetImpl is the implement of ModifySet
type ModifySetImpl struct {
	SeqNo uint64                     `json:"seq_no"`
	Data  map[string]*ModifyItemImpl `json:"data"`
}

// NewModifySet generate a new ModifySet (mainly for stateDB)
func NewModifySet() ModifySet {
	return &ModifySetImpl{
		Data: make(map[string]*ModifyItemImpl),
	}
}

// PersistToFile will persist ModifySetImpl into target file
// First line is hash of json string.
// Second line is json string of ModifySetImpl.
func (m *ModifySetImpl) PersistToFile(path string) (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bs, _ := json.Marshal(m)
	h := hash(bs)
	var f *os.File
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}

	defer func() {
		_ = f.Close()
	}()

	toWrite := bytes.NewBuffer(h)
	toWrite.Write([]byte("\n"))
	toWrite.Write(bs)
	_, err = f.Write(toWrite.Bytes())
	return err
}

func hash(bs []byte) []byte {
	hasher := sha1.New()
	_, _ = hasher.Write(bs)
	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}

// RecoverFromFile will recover ModifySet from target file
// notice that origin data will be reset.
func (m *ModifySetImpl) RecoverFromFile(path string) (err error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	bs, _ := ioutil.ReadFile(path)
	bss := bytes.Split(bs, []byte("\n"))
	if len(bss) != 2 {
		return fmt.Errorf("wrong file format")
	}

	h, ms := bss[0], bss[1]
	if !bytes.Equal(h, hash(ms)) {
		return fmt.Errorf("wrong file content")
	}

	err = json.Unmarshal(ms, m)
	if err != nil {
		return err
	}
	return nil
}

// Set set the modify items with given key, old value and new value
func (m *ModifySetImpl) Set(key string, oldVal, newVal []byte) {
	if item, exist := m.Data[key]; exist {
		item.SetNewVal(newVal)
		item.SetOldVal(oldVal)
		return
	}
	m.Data[key] = &ModifyItemImpl{
		Name:   key,
		OldVal: oldVal,
		NewVal: newVal,
	}
	return

}

// Get returns the modify items with given key
func (m *ModifySetImpl) Get(key string) (mi ModifyItem, exist bool) {
	mi, exist = m.Data[key]
	return
}

// GetAll returns the all modify items of modify set
func (m *ModifySetImpl) GetAll() map[string]ModifyItem {
	newMap := make(map[string]ModifyItem)
	for k, v := range m.Data {
		newMap[k] = v
	}
	return newMap
}

// GetSeqNo returns the sequence number of modify set
func (m *ModifySetImpl) GetSeqNo() uint64 {
	return m.SeqNo
}

// SetSeqNo sets the sequence number of modify set
func (m *ModifySetImpl) SetSeqNo(seqNo uint64) {
	m.SeqNo = seqNo
}

// Len returns the length of modify set
func (m *ModifySetImpl) Len() int {
	return len(m.Data)
}

// Del deletes certain key of modify set
func (m *ModifySetImpl) Del(key string) {
	delete(m.Data, key)
}

// Prune removes all modify items with the same old value and new value
func (m *ModifySetImpl) Prune() {
	for k, v := range m.Data {
		if bytes.Equal(v.GetOldVal(), v.GetNewVal()) {
			delete(m.Data, k)
		}
	}
}

// ModifyItemImpl is the implement of ModifyItem
type ModifyItemImpl struct {
	Name   string `json:"name"`
	OldVal []byte `json:"old"`
	NewVal []byte `json:"new"`
}

// SetNewVal sets the new value of ModifyItem
func (m *ModifyItemImpl) SetNewVal(newVal []byte) {
	m.NewVal = newVal
}

// SetOldVal sets the old value of ModifyItem
func (m *ModifyItemImpl) SetOldVal(oldVal []byte) {
	m.OldVal = oldVal
}

// GetNewVal returns the new value of ModifyItem
func (m *ModifyItemImpl) GetNewVal() []byte {
	return m.NewVal
}

// GetOldVal returns the old value of ModifyItem
func (m *ModifyItemImpl) GetOldVal() []byte {
	return m.OldVal
}

// GetName returns the name of ModifyItem
func (m *ModifyItemImpl) GetName() string {
	return m.Name
}

// ResetMetaFile clear the content of meta file
func ResetMetaFile(path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	_ = f.Close()
}
