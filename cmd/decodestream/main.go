package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/spf13/pflag"
	"github.com/yujuncen/km/pkg/keyutils"
)

type KVStream struct {
	read io.Reader
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
	buf := [8]byte{}
	_, err := s.read.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return int(binary.LittleEndian.Uint64(buf[:])), nil
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

func RefineValue(value []byte) string {
	if len(value) > 32 {
		return fmt.Sprintf("%s...(and %d bytes more)", hex.EncodeToString(value[:32]), len(value)-32)
	}
	return hex.EncodeToString(value)
}

var (
	file   = pflag.String("file", "", "The file to parse.")
	decode = pflag.Bool("decode", true, "Decode the table key.")
)

func main() {
	pflag.Parse()

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}
	kvs := KVStream{read: f}
	key := bytes.NewBuffer(nil)
	value := bytes.NewBuffer(nil)
	for {
		err := kvs.ReadKeyValue(key, value)
		if err != nil {
			if err != io.EOF {
				fmt.Println("err: ", err)
			}
			return
		}
		bKey := key.Bytes()
		var hk keyutils.HumanKey
		if *decode {
			hk = keyutils.ParseKeyFromEncodedWithTS(bKey)
		} else {
			hk = keyutils.ParseKeyForHuman(bKey)
		}
		fmt.Printf("%s => %s\n", hk.String(), RefineValue(value.Bytes()))
	}
}
