package xunit

import (
	"errors"
	"fmt"
	"math"
)

var _ = fmt.Print

type DataSize uint64
type SISize DataSize
type IECSize DataSize

const (
	// Decimal
	B = 1

	KB SISize = 1000 * SISize(B)
	MB        = 1000 * KB
	GB        = 1000 * MB
	TB        = 1000 * GB
	PB        = 1000 * TB
	EB        = 1000 * PB

	// Binary
	KiB IECSize = 1024 * IECSize(B)
	MiB         = 1024 * KiB
	GiB         = 1024 * MiB
	TiB         = 1024 * GiB
	PiB         = 1024 * TiB
	EiB         = 1024 * PiB
)

var iecM = map[IECSize]string{IECSize(B): `B`, KiB: `KiB`, MiB: `MiB`, GiB: `GiB`, TiB: `TiB`, PiB: `PiB`, EiB: `EiB`}
var siStr = map[SISize]string{SISize(B): `B`, KB: `kB`, MB: `MB`, GB: `GB`, TB: `TB`, PB: `PB`, EB: `EB`}
var strIEC = map[string]IECSize{`B`: IECSize(B), `KiB`: KiB, `MiB`: MiB, `GiB`: GiB, `TiB`: TiB, `PiB`: PiB, `EiB`: EiB}
var strSI = map[string]SISize{`B`: SISize(B), `kB`: KB, `MB`: MB, `GB`: GB, `TB`: TB, `PB`: PB, `EB`: EB}

var ErrInvalidUnit = errors.New(`unsupported unit`)
var ErrInvalidSIString = errors.New(`invalid SI string`)
var ErrInvalidIECString = errors.New(`invalid IEC string`)
var ErrOverflow = errors.New(`overflow`)

func (s SISize) String() string {
	var str string
	switch {
	case s >= EB:
		str, _ = s.Format(EB, -1, 0)
	case s >= PB:
		str, _ = s.Format(PB, -1, 0)
	case s >= TB:
		str, _ = s.Format(TB, -1, 0)
	case s >= GB:
		str, _ = s.Format(GB, -1, 0)
	case s >= MB:
		str, _ = s.Format(MB, -1, 0)
	case s >= KB:
		str, _ = s.Format(KB, -1, 0)
	default:
		str, _ = s.Format(1, -1, 0)
	}
	return str
}

// Format convert SISize to string. prec 0 means no fractional part. prec 1 means
// one fractional part digit, and so on. prec -1 means remain all of fractional part digits.
func (s SISize) Format(unit SISize, prec int, seg byte) (string, error) {
	var n int
	var low uint64
	var arr [32]byte
	arr[31] = 'B'
	switch prec {
	case -1:
		low = 1
	case 0:
		low = uint64(unit)
	default:
		low = uint64(unit) / uint64(math.Pow10(int(prec)))
	}
	switch unit {
	case EB:
		arr[30] = 'E'
		_, p := fmtFrac2(arr[:30], uint64(s%unit), uint64(unit), low)
		n = fmtUint(arr[:p], uint64(s/unit), seg)
	case PB:
		arr[30] = 'P'
		_, p := fmtFrac2(arr[:30], uint64(s%unit), uint64(unit), low)
		n = fmtUint(arr[:p], uint64(s/unit), seg)
	case TB:
		arr[30] = 'T'
		_, p := fmtFrac2(arr[:30], uint64(s%unit), uint64(unit), low)
		n = fmtUint(arr[:p], uint64(s/unit), seg)
	case GB:
		arr[30] = 'G'
		_, p := fmtFrac2(arr[:30], uint64(s%unit), uint64(unit), low)
		n = fmtUint(arr[:p], uint64(s/unit), seg)
	case MB:
		arr[30] = 'M'
		_, p := fmtFrac2(arr[:30], uint64(s%unit), uint64(unit), low)
		n = fmtUint(arr[:p], uint64(s/unit), seg)
	case KB:
		arr[30] = 'k'
		_, p := fmtFrac2(arr[:30], uint64(s%unit), uint64(unit), low)
		n = fmtUint(arr[:p], uint64(s/unit), seg)
	case SISize(B):
		n = fmtUint(arr[:31], uint64(s), seg)
	default:
		return ``, ErrInvalidUnit
	}
	return string(arr[n:]), nil
}

func (s IECSize) String() string {
	switch {
	case s >= EiB:
		return fmt.Sprintf("%.4g%s", float64(s)/float64(EiB), `EiB`)
	case s >= PiB:
		return fmt.Sprintf("%.4g%s", float64(s)/float64(PiB), `PiB`)
	case s >= TiB:
		return fmt.Sprintf("%.4g%s", float64(s)/float64(TiB), `TiB`)
	case s >= GiB:
		return fmt.Sprintf("%.4g%s", float64(s)/float64(GiB), `GiB`)
	case s >= MiB:
		return fmt.Sprintf("%.4g%s", float64(s)/float64(MiB), `MiB`)
	case s >= KiB:
		return fmt.Sprintf("%.4g%s", float64(s)/float64(KiB), `KiB`)
	default:
		return fmt.Sprintf("%d%s", uint64(s), `B`)
	}
}

func MustParseSI(s string) SISize {
	si, err := ParseSI(s)
	if err != nil {
		panic(err)
	}
	return si
}

func ParseSI(s string) (SISize, error) {
	if len(s) == 0 {
		return 0, nil
	}
	var n uint64
	var dot int
	var sep, unitprefix, digit bool
	for i := range s {
		switch s[i] {
		case ' ': //ignore space
		case ',', '_':
			if sep || n == 0 || dot > 0 {
				return 0, ErrInvalidSIString
			}
			sep = true
			continue
		case '.':
			if sep {
				return 0, ErrInvalidSIString
			}
			dot += 1
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			sep = false
			n *= 10
			n += uint64(s[i] - '0')
			if dot != 0 {
				dot += 1
			}
			digit = true
		case '0':
			if n != 0 {
				n *= 10
			}
			if dot > 0 {
				dot += 1
			}
			digit = true
		case 'k', 'M', 'G', 'T', 'P', 'E':
			if n == 0 && !digit {
				return 0, ErrInvalidSIString
			}
			if dot == 1 {
				return 0, ErrInvalidSIString
			}
			if i == len(s) || s[i+1] != 'B' {
				return 0, ErrInvalidSIString
			}
			a, b := uint64(strSI[string(s[i:i+2])]), uint64(math.Pow10(max(dot-1, 0)))
			if a > b {
				n *= a / b
			} else {
				n /= b / a
			}
			unitprefix = true
		case 'B':
			if i < 1 || (!unitprefix && dot > 0) || sep {
				return 0, ErrInvalidSIString
			}
		default:
			return 0, ErrInvalidSIString
		}
	}
	return SISize(n), nil
}

func fmtUint(buf []byte, n uint64, s byte) int {
	var d uint64
	var p = len(buf)
	var c int
NEXT:
	n, d = n/10, n%10
	p -= 1
	c += 1
	buf[p] = '0' + byte(d)
	if c%3 == 0 && s != 0 {
		p -= 1
		buf[p] = s
	}
	if n != 0 {
		goto NEXT
	}
	return p
}

type formatOpt struct {
	u SISize
	s byte
	p uint8
}
type FormatOption func(o *formatOpt)

func WithSIUnit(unit SISize) FormatOption {
	return func(o *formatOpt) { o.u = unit }
}
func WithSegment(s byte) FormatOption {
	return func(o *formatOpt) { o.s = s }
}
func WithPrecison(p uint8) FormatOption {
	return func(o *formatOpt) { o.p = p }
}

// 1.00 MB = 1,003,000 B, digit = 2
// n = 3000, high = 1,000,000 ,low = 10,000
// buf should be '00'
//
// 1 MB = 1,003,000 B, digit = 0
// n = 3000, high = 1,000,000 ,low = 1,000,000
//
// 1.0 MB = 1,003,000 B, digit = 1
// n = 3000, high = 1,000,100 ,low = 100,000
//
// 1.003 MB = 1,003,000 B, digit = -1
// n = 3000, high = 1,000,000 ,low = 1
func fmtFrac2(buf []byte, n uint64, high uint64, low uint64) (uint64, int) {
	var (
		p     = len(buf)
		d     uint64
		print bool = low != 1
	)
	high /= low
	n /= low
	for i := uint64(1); i < high; i *= 10 {
		n, d = n/10, n%10
		if !print && d == 0 {
			continue
		}
		print = true
		p -= 1
		buf[p] = '0' + byte(d)
	}
	if print && high != 1 {
		p -= 1
		buf[p] = '.'
	}
	return n, p
}

func MustParseIEC(s string) IECSize {
	iec, err := ParseIEC(s)
	if err != nil {
		panic(err)
	}
	return iec
}

func ParseIEC(s string) (IECSize, error) {
	if len(s) == 0 {
		return 0, nil
	}
	var n uint64
	var dot int
	var sep, unitprefix, digit bool
	for i := range s {
		switch s[i] {
		case ' ': //ignore space
		case ',', '_':
			if sep || n == 0 || dot > 0 {
				return 0, ErrInvalidIECString
			}
			sep = true
			continue
		case '.':
			if sep {
				return 0, ErrInvalidIECString
			}
			dot += 1
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			sep = false
			digit = true
			n *= 10
			n += uint64(s[i] - '0')
			if dot != 0 {
				dot += 1
			}
		case '0':
			if n != 0 {
				n *= 10
			}
			if dot > 0 {
				dot += 1
			}
			digit = true
		case 'K', 'M', 'G', 'T', 'P', 'E':
			if n == 0 && !digit {
				return 0, ErrInvalidIECString
			}
			if dot == 1 {
				return 0, ErrInvalidIECString
			}
			if i >= len(s)+1 || s[i+1] != 'i' || s[i+2] != 'B' {
				return 0, ErrInvalidIECString
			}
			unit := strIEC[string(s[i:i+3])]
			n = uint64(float64(n) / math.Pow10(max(dot-1, 0)) * float64(unit))
			unitprefix = true
		case 'i':
			if i == len(s) || s[i+1] != 'B' {
				return 0, ErrInvalidIECString
			}
		case 'B':
			if i < 1 || (!unitprefix && dot > 0) || sep {
				return 0, ErrInvalidIECString
			}
		default:
			return 0, ErrInvalidIECString
		}
	}
	return IECSize(n), nil
}
