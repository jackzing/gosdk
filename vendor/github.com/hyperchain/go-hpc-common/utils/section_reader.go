// Hyperchain License
// Copyright (C) 2016 The Hyperchain Authors.

package utils

import (
	"errors"
	"io"
	"os"
)

// SectionReader is used to read the shard of a file
type SectionReader struct {
	filePath    string
	shardLen    int64
	latestShard int64
	shardNum    int64
	fd          *os.File
}

// Err type
var (
	ErrInvalidShardID = errors.New("invalid shard id")
	ErrEmptyFile      = errors.New("empty file")
	ErrNotFile        = errors.New("is not file")
)

// NewSectionReader returns a SectionReader, contains filePath and an open file.
func NewSectionReader(filePath string, shardLen int64) (*SectionReader, error) {
	fd, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	fstat, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	if fstat.IsDir() {
		return nil, ErrNotFile
	}

	size := fstat.Size()
	if size == 0 {
		return nil, ErrEmptyFile
	}

	shardNum := size / shardLen
	if (size % shardLen) > 0 {
		shardNum++
	}

	return &SectionReader{
		filePath: filePath,
		shardLen: shardLen,
		fd:       fd,
		shardNum: shardNum,
	}, nil
}

// ReadNext reads the next shard of the latestShard
func (sectionReader *SectionReader) ReadNext() (n int, buf []byte, err error) {
	defer func() {
		if err == nil {
			sectionReader.latestShard++
		}
	}()

	if !sectionReader.check(sectionReader.latestShard + 1) {
		n = 0
		buf = nil
		err = ErrInvalidShardID
		return
	}

	buf = make([]byte, sectionReader.shardLen)
	offset := sectionReader.latestShard * sectionReader.shardLen
	r := io.NewSectionReader(sectionReader.fd, offset, sectionReader.shardLen)
	n, err = r.Read(buf)
	return
}

// ReadAt reads the sidth shard
func (sectionReader *SectionReader) ReadAt(sid int64) (int, []byte, error) {
	sectionReader.latestShard = sid - 1
	return sectionReader.ReadNext()
}

// Close closes the open file
func (sectionReader *SectionReader) Close() {
	// nolint
	if sectionReader.fd != nil {
		sectionReader.fd.Close()
	}
}

// check checks that if sid is legal or not
func (sectionReader *SectionReader) check(sid int64) bool {
	if sid == 0 || sid > sectionReader.shardNum {
		return false
	}
	return true
}
