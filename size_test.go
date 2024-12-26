package xunit_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/violin0622/xunit"
)

var _ = fmt.Println

func TestSIString(t *testing.T) {
	cases := []struct {
		size xunit.SISize
		str  string
	}{
		{xunit.SISize(0), `0B`},
		{xunit.SISize(xunit.B), `1B`},
		{xunit.SISize(17 * xunit.B), `17B`},
		{xunit.KB, `1kB`},
		{115 * xunit.MB, `115MB`},
		{1000 * xunit.MB, `1GB`},
		{xunit.GB, `1GB`},
		{xunit.TB, `1TB`},
		{xunit.EB, `1EB`},
		{xunit.SISize(math.MaxUint64), `18.446744073709551615EB`},
		{xunit.MB + xunit.KB, `1.001MB`},
		{xunit.SISize(1.5 * float64(xunit.MB)), `1.5MB`},
	}

	assert := asserter{t}
	for _, c := range cases {
		assert.equal(c.size.String(), c.str)
	}
}

func BenchmarkSIString(b *testing.B) {
	var si xunit.SISize = 123 * xunit.TB
	for i := 0; i < b.N; i++ {
		if si.String() != `123TB` {
			b.Fail()
		}
	}
}

func TestIECString(t *testing.T) {
	cases := []struct {
		size xunit.IECSize
		str  string
	}{
		{xunit.IECSize(0), `0B`},
		{xunit.IECSize(xunit.B), `1B`},
		{xunit.IECSize(17 * xunit.B), `17B`},
		{xunit.KiB, `1KiB`},
		{115 * xunit.MiB, `115MiB`},
		{1024 * xunit.MiB, `1GiB`},
		{xunit.GiB, `1GiB`},
		{xunit.TiB, `1TiB`},
		{xunit.EiB, `1EiB`},
		{xunit.IECSize(math.MaxUint64), `16EiB`},
		{xunit.MiB + xunit.KiB, `1.001MiB`},
		{xunit.MiB + xunit.IECSize(1), `1MiB`},
		{xunit.IECSize(1.5 * float64(xunit.MiB)), `1.5MiB`},
	}

	assert := asserter{t}
	for _, c := range cases {
		assert.equal(c.size.String(), c.str)
	}
}

func BenchmarkIECString(b *testing.B) {
	var iec xunit.IECSize = xunit.IECSize(123.5 * float64(xunit.TiB))
	for i := 0; i < b.N; i++ {
		if iec.String() != `123.5TiB` {
			b.Fail()
		}
	}
}

func TestParseSI(t *testing.T) {
	cases := []struct {
		size xunit.SISize
		str  string
	}{

		{xunit.SISize(0), ``},
		{xunit.SISize(0), `0B`},
		{xunit.SISize(0), `0kB`},
		{xunit.SISize(500), `0.5kB`},
		{xunit.SISize(123), `123 B`},
		{xunit.SISize(1230), ` 1,2 30 B `},
		{xunit.SISize(1230), ` 1.23 kB `},
		{xunit.SISize(1230), ` 1.2300 kB `},
		{xunit.SISize(11230), ` 11.2300 kB `},
		{xunit.SISize(xunit.B), `1B`},
		{xunit.SISize(17 * xunit.B), `17B`},
		{xunit.KB, `1kB`},
		{115 * xunit.MB, `115MB`},
		{1000 * xunit.MB, `1GB`},
		{xunit.GB, `1GB`},
		{xunit.TB, `1TB`},
		{xunit.EB, `1EB`},
		{xunit.SISize(math.MaxUint64), `18.446744073709551615EB`},
		{xunit.MB + xunit.KB, `1.001MB`},
		{xunit.SISize(1.5 * float64(xunit.MB)), `1.5MB`},
	}
	assert := asserter{t}
	for i, c := range cases {
		size, e := xunit.ParseSI(c.str)
		assert.nil(e, `%d '%s'`, i, c.str)
		assert.equal(size, c.size, `%d '%s'`, i, c.str)
	}
}

func BenchmarkParseSI(b *testing.B) {
	for i := 0; i < b.N; i++ {
		size, err := xunit.ParseSI(` 11.2300 kB `)
		if err != nil || size != 11230 {
			b.Fail()
		}
	}
}

func BenchmarkParseIEC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		size, err := xunit.ParseIEC(` 1.1 KiB `)
		if err != nil || size != 1126 {
			b.Fail()
		}
	}
}

func TestParseSIInvalid(t *testing.T) {
	// theses string is invalid.
	cases := []string{
		`,`,
		`B`,
		`kB`,
		`1.2B`,
		`,33B`,
		`33,B`,
		`3,,3B`,
		`33.3k B`,
		`33.3KB`,
		`33.3k_B`,
		`33.3k,B`,
		`33,.3kB`,
		`33,.,3kB`,
		`33,.,3kB`,
		`33.33,3kB`,
		`_`,
		`_33B`,
		`33_B`,
		`3__3B`,
	}
	assert := asserter{t}
	var i int
	var c string
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(i, c)
		}
	})
	for i, c = range cases {
		size, e := xunit.ParseSI(c)
		assert.equal(e, xunit.ErrInvalidSIString, `%d '%s'`, i, c)
		assert.equal(size, xunit.SISize(0))
	}
}

func TestParseIEC(t *testing.T) {
	cases := []struct {
		size xunit.IECSize
		str  string
	}{
		{xunit.IECSize(0), ``},
		{xunit.IECSize(0), `0B`},
		{xunit.IECSize(0), `0KiB`},
		{xunit.IECSize(512), `0.5KiB`},
		{xunit.IECSize(123), `123 B`},
		{xunit.IECSize(1230), ` 1,2 30 B `},
		{xunit.IECSize(1280), ` 1.25 KiB `},
		{xunit.IECSize(1280), ` 1.2500 KiB `},
		{xunit.IECSize(11520), ` 11.2500 KiB `},
		{xunit.KiB, `1KiB`},
		{115 * xunit.MiB, `115MiB`},
		{1024 * xunit.MiB, `1GiB`},
		{xunit.GiB, `1GiB`},
		{xunit.TiB, `1TiB`},
		{xunit.EiB, `1EiB`},
		{xunit.IECSize(math.MaxUint64), `16EiB`},     //16EiB is the maximum supported IEC size.
		{xunit.IECSize(math.MaxUint64), `16,384PiB`}, //16EiB is the maximum supported IEC size.
		{xunit.MiB + xunit.KiB + 24*xunit.IECSize(xunit.B), `1.001MiB`},
		{xunit.MiB + 1023*xunit.MiB, `1GiB`},
		{xunit.IECSize(1.5 * float64(xunit.MiB)), `1.5MiB`},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf(`%d_%s_%v`, i, c.str, c.size), func(t *testing.T) {
			assert := asserter{t}
			size, e := xunit.ParseIEC(c.str)
			assert.nil(e)
			assert.equal(size, c.size)
		})
	}
}

func TestParseIECInvalid(t *testing.T) {
	// theses string is invalid.
	cases := []string{
		`,`,
		`B`,
		`kB`,
		`KB`,
		`KiB`,
		`1.2B`,
		`,33B`,
		`33,B`,
		`3,,3B`,
		`33.3Ki B`,
		`33.3K iB`,
		`33.3Ki_B`,
		`33.3K_iB`,
		`33.3K,iB`,
		`33.3Ki,B`,
		`33,.3KiB`,
		`33,.,3KiB`,
		`33,.,3KiB`,
		`33.33,3KiB`,
		`_`,
		`_33B`,
		`33_B`,
		`3__3B`,
	}
	var i int
	var c string
	t.Cleanup(func() {
		if t.Failed() {
			t.Log(i, c)
		}
	})
	for i, c = range cases {
		t.Run(fmt.Sprintf(`%d_%s`, i, c), func(t *testing.T) {
			assert := asserter{t}
			size, e := xunit.ParseIEC(c)
			assert.equal(e, xunit.ErrInvalidIECString)
			assert.equal(size, xunit.IECSize(0))
		})
	}
}

func TestFormatSI(t *testing.T) {
	cases := []struct {
		size      xunit.SISize
		str       string
		unit      xunit.SISize
		precision int
		seg       byte
	}{
		{xunit.KB, `1kB`, xunit.KB, 0, 0},
		{xunit.KB, `1kB`, xunit.KB, 0, ','},
		{xunit.KB, `1000B`, xunit.SISize(xunit.B), 0, 0},
		{xunit.KB, `1,000B`, xunit.SISize(xunit.B), 0, ','},
		{xunit.MB, `1MB`, xunit.MB, 0, 0},
		{1204 * xunit.MB, `1204MB`, xunit.MB, 0, 0},
		{1204 * xunit.MB, `1204000000B`, xunit.SISize(xunit.B), 0, 0},
		{1204 * xunit.MB, `1,204,000,000B`, xunit.SISize(xunit.B), 0, ','},
		{1204 * xunit.MB, `1_204_000_000B`, xunit.SISize(xunit.B), 0, '_'},
		{xunit.MB + 3*xunit.KB, `1.003MB`, xunit.MB, -1, 0},
		{xunit.MB + 3*xunit.KB, `1.00MB`, xunit.MB, 2, 0},
		{xunit.MB + 3*xunit.KB, `1.0MB`, xunit.MB, 1, 0},
		{xunit.MB + 3*xunit.KB, `1MB`, xunit.MB, 0, 0},
		{xunit.MB + xunit.KB, `1001kB`, xunit.KB, 0, 0},
		{xunit.MB + xunit.KB + xunit.SISize(1), `1.001001MB`, xunit.MB, -1, 0},
		{xunit.MB + 12*xunit.KB, `1.012MB`, xunit.MB, -1, 0},
		{xunit.SISize(1.5 * float64(xunit.MB)), `1.5MB`, xunit.MB, -1, 0},
		{xunit.SISize(1.5 * float64(xunit.MB)), `1500kB`, xunit.KB, -1, 0},
	}

	assert := asserter{t}
	for _, c := range cases {
		a, e := c.size.FormatSI(c.unit, c.precision, c.seg)
		assert.nil(e)
		assert.equal(a, c.str)
	}
}

type asserter struct {
	t *testing.T
}

func (a asserter) nil(actual any, ex ...any) {
	a.t.Helper()
	if actual != nil {
		switch len(ex) {
		case 0:
			a.t.Errorf(`Not nil. Actual: %v`, actual)
		case 1:
			a.t.Errorf("%v\nNot nil. Actual: %v", ex[0], actual)
		default:
			a.t.Errorf("%v\nNot nil. Actual: %v", fmt.Sprintf(ex[0].(string), ex[1:]...), actual)
		}
	}
}

func (a asserter) equal(actual, expect any, ex ...any) {
	a.t.Helper()
	if actual != expect {
		switch len(ex) {
		case 0:
			a.t.Errorf(`Expect: %v, Actual: %v`, expect, actual)
		case 1:
			a.t.Errorf("%v\nExpect: %v, Actual: %v, ", ex[0], expect, actual)
		default:
			a.t.Errorf("%s\nExpect: %v, Actual: %v, ", fmt.Sprintf(ex[0].(string), ex[1:]...), expect, actual)
		}
	}
}

func assertEqual(t *testing.T) {
	t.Helper()

}
