package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/tikv/pd/pkg/codec"
	"github.com/yujuncen/km/pkg/keyutils"
)

var (
	cmd       = pflag.String("cmd", "decode", "the command for use.")
	file      = pflag.String("file", "", "The file to parse.")
	decode    = pflag.Bool("decode", true, "Decode the table key.")
	useJSON   = pflag.Bool("json", false, "Use Json format.")
	searchKey = pflag.BytesHex("search", nil, "The key for search. (used for `search` command.)")
)

func RefineValue(value []byte) string {
	if len(value) > 32 {
		return fmt.Sprintf("%s...(and %d bytes more)", hex.EncodeToString(value[:32]), len(value)-32)
	}
	return hex.EncodeToString(value)
}

func decodeFile(kvs keyutils.KVStream) {
	kvs.Iterate(func(key, value *bytes.Buffer) {
		bKey := key.Bytes()
		var hk keyutils.HumanKey
		if *decode {
			hk = keyutils.ParseKeyFromEncodedWithTS(bKey)
		} else {
			hk = keyutils.ParseKeyForHuman(bKey)
		}
		if *useJSON {
			b, err := json.Marshal(hk)
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		} else {
			fmt.Printf("%s => %s\n", hk.String(), RefineValue(value.Bytes()))
		}
	})
}

func search(db string) error {
	target := *searchKey
	if *decode {
		target = codec.EncodeBytes(target)
	}
	return filepath.WalkDir(db, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d == nil {
			return errors.New("nil entry: " + path)
		}
		if d.Type().IsRegular() && strings.HasSuffix(d.Name(), ".log") {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			kvs := keyutils.NewStreamReader(f)
			kvs.Iterate(func(key, value *bytes.Buffer) {
				if bytes.Equal(key.Bytes()[:key.Len()-8], target) {
					hk := keyutils.ParseKeyFromEncodedWithTS(key.Bytes())
					fmt.Printf("Key: %s\nValue: %s", hk.String(), hex.EncodeToString(value.Bytes()))
				}
			})
		}
		return nil
	})
}

func main() {
	pflag.Parse()

	switch *cmd {
	case "decode":
		f, err := os.Open(*file)
		if err != nil {
			panic(err)
		}
		decodeFile(keyutils.NewStreamReader(f))
	case "search":
		if err := search(*file); err != nil {
			panic(err)
		}
	default:
		panic("failed to decode the command.")
	}
}
