package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pingcap/log"
	"github.com/tikv/pd/pkg/codec"
	"go.uber.org/zap"
)

func Lines(f func(string)) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		f(s.Text())
	}
}

func LineHexes(f func([]byte)) {
	Lines(func(s string) {
		hexKey, err := hex.DecodeString(s)
		if err != nil {
			log.Warn("cannot decode hex key", zap.Error(err))
			return
		}
		f(hexKey)
	})
}

func main() {
	LineHexes(func(b []byte) {
		data := struct {
			Decoded string `json:"decoded"`
			Remain  string `json:"remained"`
		}{}
		remain, decoded, err := codec.DecodeBytes(b)
		if err != nil {
			log.Fatal("error during decode", zap.Error(err))
		}
		data.Decoded = fmt.Sprintf("%X", decoded)
		data.Remain = fmt.Sprintf("%X", remain)
		j, err := json.Marshal(data)
		if err != nil {
			log.Warn("cannot marshal json", zap.Error(err))
			return
		}
		fmt.Println(string(j))
	})
}
