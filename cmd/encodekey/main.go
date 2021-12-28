package main

import (
	"bufio"
	"encoding/hex"
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
		fmt.Printf("%X\n", codec.EncodeBytes(b))
	})
}
