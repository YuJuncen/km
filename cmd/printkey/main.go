package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pingcap/log"
	"github.com/yujuncen/km/pkg/keyutils"
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
	LineHexes(func(hexKey []byte) {
		key := keyutils.ParseKeyForHuman(hexKey)
		j, err := json.Marshal(key)
		if err != nil {
			log.Warn("cannot encode json", zap.Error(err))
		}
		fmt.Println(string(j))
	})
}
