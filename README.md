# Introduction

KM(Key Manager) is a tool for decoding TiDB format keys.

## Usage

### decode-stream

The tool is for decoding the KV stream format generated by stream backup.

build (Aha, no Makefiles yet.):

```bash
build -o decode-stream cmd/decodestream/main.go
```

usage: 

```bash
> decode-stream --file 00000274_00000920_write_Put.temp.log | head
Index(table=274, index=1, index_data=[9 1])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000051
Index(table=274, index=1, index_data=[9 2])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000052
Index(table=274, index=1, index_data=[9 3])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000053
Index(table=274, index=1, index_data=[9 4])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000054
Index(table=274, index=1, index_data=[9 5])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000055
Index(table=274, index=1, index_data=[9 6])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000056
Index(table=274, index=1, index_data=[9 7])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000057
Index(table=274, index=1, index_data=[9 8])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000058
Index(table=274, index=1, index_data=[9 9])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb0576080000000000000059
Index(table=274, index=1, index_data=[9 10])@(phy=2021-12-28 01:03:19,log=3) => 508280e0dcacecfcfb057608000000000000005a
```

### (Doc of other tools TBD)