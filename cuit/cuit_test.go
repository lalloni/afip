// MIT License
//
// Copyright (c) 2018 Pablo Ignacio Lalloni
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cuit

import (
	"fmt"
	"math/rand"
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
)

func ExampleIsValid_invalid() {
	v := IsValid(20123456781)
	fmt.Println(v)
	// Output: false
}

func ExampleIsValid_valid() {
	v := IsValid(33693450239)
	fmt.Println(v)
	// Output: true
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name string
		cuit uint64
		want bool
	}{
		{"valid", 20242643772, true},
		{"bad verifier", 20242643773, false},
		{"valid woman", 27240366180, true},
		{"bad verifier woman", 27240366185, false},
		{"valid legal", 33693450239, true},
		{"invalid legal", 33603450239, false},
		{"valid legal 2", 30711413568, true},
		{"invalid legal", 31711413568, false},
		{"invalid range big", 10030711413568, false},
		{"invalid range small", 100, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, IsValid(test.cuit))
		})
	}
}

func ExampleParts_basic() {
	kind, id, ver := Parts(20123456781)
	fmt.Println(kind, id, ver)
	// Output: 20 12345678 1
}

func ExampleParts_range() {
	kind, id, ver := Parts(10020123456781)
	fmt.Println(kind, id, ver)
	// Output: 20 12345678 1
}

func TestParts(t *testing.T) {
	tests := []struct {
		name string
		cuit uint64
		kind uint64
		id   uint64
		ver  uint64
	}{
		{"basic", 10123456781, 10, 12345678, 1},
		{"zero verifier", 10123456780, 10, 12345678, 0},
		{"small padded id", 20003456782, 20, 345678, 2},
		{"discard excess digits", 10027876543215, 27, 87654321, 5},
		{"big zero", 1e18, 0, 0, 0},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			kind, id, ver := Parts(test.cuit)
			assert.Equal(t, test.kind, kind)
			assert.Equal(t, test.id, id)
			assert.Equal(t, test.ver, ver)
		})
	}
	f := func(cuit uint64) bool {
		kind, id, ver := Parts(cuit)
		c := cuit % 1e11
		return assert.Equal(t, c/1e9, kind) &&
			assert.Equal(t, (c%1e9)/10, id) &&
			assert.Equal(t, c%10, ver)
	}
	t.Run("quickchecks", func(t *testing.T) {
		if err := quick.Check(f, &quick.Config{MaxCount: 1000}); err != nil {
			t.Error(err)
		}
	})
}

func ExampleParse() {
	cuit, err := Parse("20-12345678-2")
	if err != nil {
		// handle parse error
	}
	fmt.Println(cuit)
	// Output: 20123456782
}

func ExampleParse_nodash() {
	cuit, err := Parse("20123456782")
	if err != nil {
		// handle parse error
	}
	fmt.Println(cuit)
	// Output: 20123456782
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		cuit    string
		want    uint64
		wantErr bool
	}{
		{"basic", "10-12345678-1", 10123456781, false},
		{"no dashes", "10123456781", 10123456781, false},
		{"one dash", "10-123456781", 10123456781, false},
		{"too big", "10010-123456781", 0, true},
		{"bad number", "1a0-12345678-x", 0, true},
		{"one", "00-00000000-1", 1, false},
		{"smallest possible", "00-00000000-0", 0, false},
		{"biggest possible", "99-99999999-9", 99999999999, false},
		{"anything", "dadaddsa", 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.cuit)
			if tt.wantErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func ExampleFormat() {
	s := Format(20123456781)
	fmt.Println(s)
	// Output: 20-12345678-1
}
func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		cuit uint64
		want string
	}{
		{"basic", 20123456781, "20-12345678-1"},
		{"zeroes", 123456781, "00-12345678-1"},
		{"padding", 10000000781, "10-00000078-1"},
		{"more padding", 1000000781, "01-00000078-1"},
		{"big", 10030123456781, "30-12345678-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Format(tt.cuit))
		})
	}
}

func ExampleVerifier() {
	fmt.Println(Verifier(201111111083123)) // bigger number
	fmt.Println(Verifier(20111111117))     // incorrect verifier
	fmt.Println(Verifier(20111111112))     // correct verifier
	fmt.Println(Verifier(20236111))        // smaller number
	// Output:
	// 8
	// 2
	// 2
	// 6
}

func ExampleRandom() {
	notrandom := rand.New(rand.NewSource(1))
	for i := 0; i < 20; i++ {
		fmt.Println(Random(notrandom))
	}
	// Output:
	// 24821535510
	// 27897858598
	// 24491673202
	// 20865519568
	// 24350552576
	// 34636692877
	// 27144588369
	// 30381828731
	// 24653598994
	// 27081291395
	// 20381499562
	// 30835156377
	// 27422311098
	// 27938316509
	// 20430747798
	// 24047844432
	// 23879884594
	// 34081375471
	// 27838147157
	// 20733877805
}

func TestRandom(t *testing.T) {
	a := assert.New(t)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000000; i++ {
		cuit := Random(rnd)
		a.True(IsValid(cuit), "invalid random cuit: %d", cuit)
	}
}

func ExampleCompose() {
	fmt.Println(Compose(20, 12345678, 1))       // basic case
	fmt.Println(Compose(20, 123456789, 1))      // excess id
	fmt.Println(Compose(210, 12345678, 1))      // excess kind
	fmt.Println(Compose(20, 12345678, 12))      // excess verifier
	fmt.Println(Compose(1234, 123456789, 1234)) // excess all
	// Output:
	// 20123456781
	// 20234567891
	// 10123456781
	// 20123456782
	// 34234567894
}
