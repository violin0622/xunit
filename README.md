# xunit
golang extention library about unit. 

## Features
- [x] (SI size) string method
- [x] (SI size) format with unit.
- [x] (SI size) format with comma.
- [x] (SI size) format with precise.
- [] (SI size) scientific notation
- [x] (SI size) parse from string
- [x] (SI size) parse from string with comma,space,underscore
- [x] (IEC size) string method
- []  (IEC size) scientific notation
- [x] (IEC size) parse from string
- [x] (IEC size) parse from string with comma,space,underscore

## Usage

Add package:  
```shell-session
go get github.com/violin0622/xunit
```

Parse string to IEC/SI size:  

```go
iecString := `12 KiB`  
iec, err := xunit.ParseIEC(iecString)
fmt.Println(iec)
// 12KiB
fmt.Println(iec == 12_288)
// true

siString := `12 kB`  
si, err := xunit.ParseSI(iecString)
fmt.Println(si)
// 12KiB
fmt.Println(si == 12_000)
// true
```

Use 'xunit' as you would 'time' package to construct an IEC/SI size,
and format it to string.

```go

iec := 150 * xunit.KiB
fmt.Println(iec.String())
// 150KiB

si := 1500 * xunit.KB
fmt.Println(si.String())
// 1.5MB

fmt.Println(si.Format(xunit.KB, 0, 0))
// 1500kB

fmt.Println(si.Format(xunit.KB, 0, ','))
// 1,500kB
```

Max supported size is `math.MaxUint64`, i.e. `16EiB` or `18.446744073709551615EB` .

## Benchmark

```shell-session
❯ go test  -bench=. -benchmem -cpu 1,2,4
goos: darwin
goarch: arm64
pkg: github.com/violin0622/xunit
cpu: Apple M1 Pro
BenchmarkSIString      	35986402	       33.48 ns/op	      5 B/op	      1 allocs/op
BenchmarkSIString-2    	35392226	       33.33 ns/op	      5 B/op	      1 allocs/op
BenchmarkSIString-4    	34977730	       33.28 ns/op	      5 B/op	      1 allocs/op
BenchmarkIECString     	12718561	       93.86 ns/op	     16 B/op	      2 allocs/op
BenchmarkIECString-2   	12967633	       92.22 ns/op	     16 B/op	      2 allocs/op
BenchmarkIECString-4   	12937447	       92.30 ns/op	     16 B/op	      2 allocs/op
BenchmarkParseSI       	41272747	       29.01 ns/op	      0 B/op	      0 allocs/op
BenchmarkParseSI-2     	41061027	       29.02 ns/op	      0 B/op	      0 allocs/op
BenchmarkParseSI-4     	41145854	       29.03 ns/op	      0 B/op	      0 allocs/op
BenchmarkParseIEC      	48399781	       24.66 ns/op	      0 B/op	      0 allocs/op
BenchmarkParseIEC-2    	48223515	       24.63 ns/op	      0 B/op	      0 allocs/op
BenchmarkParseIEC-4    	48732696	       24.69 ns/op	      0 B/op	      0 allocs/op
```
