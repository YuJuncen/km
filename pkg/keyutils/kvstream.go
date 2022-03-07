package keyutils

import (
	"bytes"
	"encoding/binary"
	"io"
)

type KVStream struct {
	read io.Reader
}

func NewStreamReader(r io.Reader) KVStream {
	return KVStream{
		read: r,
	}
}

func (s KVStream) ReadSize() (int, error) {
	buf := [4]byte{}
	_, err := s.read.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint32(buf[:])), nil
}

func (s KVStream) ReadValueSize() (int, error) {
	buf := [4]byte{}
	_, err := s.read.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint32(buf[:])), nil
}

func (s KVStream) ReadChunk(chunk *bytes.Buffer, keyLen int) error {
	chunk.Reset()
	chunk.Grow(keyLen)
	_, err := chunk.ReadFrom(io.LimitReader(s.read, int64(keyLen)))
	return err
}

func (s KVStream) ReadKeyValue(key, value *bytes.Buffer) error {
	keyLen, err := s.ReadSize()
	if err != nil {
		return err
	}
	if err := s.ReadChunk(key, keyLen); err != nil {
		return err
	}
	valueLen, err := s.ReadValueSize()
	if err != nil {
		return err
	}
	if err := s.ReadChunk(value, valueLen); err != nil {
		return err
	}
	return nil
}

func (s KVStream) Iterate(i func(key, value *bytes.Buffer)) error {
	key := bytes.NewBuffer(nil)
	value := bytes.NewBuffer(nil)
	for {
		err := s.ReadKeyValue(key, value)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}

		i(key, value)
	}
}
