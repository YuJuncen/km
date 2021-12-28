package keyutils

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/pingcap/tidb/tablecodec"
	"github.com/tikv/pd/pkg/codec"
	"github.com/tikv/pd/pkg/tsoutil"
)

type HumanKey struct {
	Type      string   `json:"type"`
	Origin    string   `json:"origin,omitempty"`
	Table     int64    `json:"table_id,omitempty"`
	Handle    string   `json:"handle,omitempty"`
	Index     int64    `json:"index_id,omitempty"`
	IndexData []string `json:"index_data,omitempty"`
	TS        uint64   `json:"ts,omitempty"`
	Encoded   bool     `json:"encoded"`
}

const (
	TypeRecord  = "record"
	TypeIndex   = "index"
	TypeUnknown = "unknown"
)

func (hk *HumanKey) String() string {
	return hk.BaseString() + hk.TsString()
}

func (hk *HumanKey) BaseString() string {
	switch hk.Type {
	case TypeRecord:
		return fmt.Sprintf("Record(table=%d, handle=%s)", hk.Table, hk.Handle)
	case TypeIndex:
		return fmt.Sprintf("Index(table=%d, index=%d, index_data=%v)", hk.Table, hk.Index, hk.IndexData)
	default:
		return fmt.Sprintf("Unknown(raw=%s)", hk.Origin)
	}
}

func (hk *HumanKey) TsString() string {
	if hk.TS == 0 {
		return ""
	}
	return "@" + HumanTS(hk.TS)
}

func HumanTS(ts uint64) string {
	physical, logic := tsoutil.ParseTS(ts)
	return fmt.Sprintf("(phy=%s,log=%d)", physical.Format("2006-01-02 15:04:05"), logic)
}

func ParseKeyFromEncodedWithTS(key []byte) HumanKey {
	d := len(key) - 8
	encoded, ts := key[:d], key[d:]
	_, decoded, err := codec.DecodeBytes(encoded)
	if err != nil {
		decoded = encoded
	}
	tsInt := ^binary.BigEndian.Uint64(ts)
	hk := ParseKeyForHuman(decoded)
	hk.TS = tsInt
	return hk
}

func ParseKeyForHuman(key []byte) HumanKey {
	origin := fmt.Sprintf("%X", key)
	if tablecodec.IsRecordKey(key) {
		tid, rid, err := tablecodec.DecodeRecordKey(key)
		if err == nil {
			return HumanKey{
				Origin: origin,
				Type:   TypeRecord,
				Table:  tid,
				Handle: rid.String(),
			}
		}
	}
	if tablecodec.IsIndexKey(key) {
		tid, iid, data, err := tablecodec.DecodeIndexKey(key)
		if err == nil {
			return HumanKey{
				Origin:    origin,
				Type:      TypeIndex,
				Table:     tid,
				Index:     iid,
				IndexData: data,
			}
		}
	}
	return HumanKey{
		Origin: origin,
		Type:   TypeUnknown,
		Table:  tablecodec.DecodeTableID(key),
	}
}

func Lines(f func(string)) {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		f(s.Text())
	}
}
