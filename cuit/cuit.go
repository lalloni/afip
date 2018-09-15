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
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	kinds   = []uint64{20, 23, 24, 27, 30, 33, 34}
	factor  = []uint64{2, 3, 4, 5, 6, 7}
	factors = len(factor)
	pattern = regexp.MustCompile(`^(\d{2})-?(\d{8})-?(\d{1})$`)
)

// IsValid checks the provided CUIT/CUIL number for validity considering
// wether the number is in valid range, has a valid kind (20, 23, 24, 27,
// 30, 33 or 34) and has a correct verifier digit.
func IsValid(cuit uint64) bool {
	return validSize(cuit) && validKind(cuit) && validVerifier(cuit)
}

func validSize(cuit uint64) bool {
	return cuit/1e11 == 0
}

func validKind(cuit uint64) bool {
	cuitk := cuit / 1e9
	valid := false
	for _, kind := range kinds {
		if cuitk == kind {
			valid = true
		}
	}
	return valid
}

func validVerifier(cuit uint64) bool {
	return cuit%10 == verifier(cuit)
}

// verifier calcula el dígito verificador que correspondería al cuit
// suministrado devolviendo un número entero de 0 a 10.
// Si devuelve un entero de 0 a 9 es el verificador correspondiente.
// Si devuelve un entero igual a 10 el cuit suministrado es incorrecto.
func verifier(cuit uint64) uint64 {
	var num uint64
	rem := cuit / 10 // drop verifier
	for i := 0; rem != 0; i++ {
		num = num + factor[i%factors]*(rem%10)
		rem = rem / 10
	}
	num = 11 - (num % 11)
	if num == 11 {
		return 0
	}
	return num
}

// Parse extracts a CUIT number from the string provided using the
// standard format "DD-DDDDDDDD-D" being "D" any decimal digit and
// both "-" characters optional.
//
// If the string can not be parsed as a CUIT number the function returns error.
func Parse(cuit string) (uint64, error) {
	match := pattern.FindStringSubmatch(cuit)
	if match == nil {
		return 0, errors.New("formato de CUIT/CUIL incorrecto")
	}
	// the following three errors can never
	// happen because the regexp pattern
	// ensures all submatches are **only**
	// digits, so we can safely ignore them
	kind, _ := strconv.Atoi(match[1])
	id, _ := strconv.Atoi(match[2])
	ver, _ := strconv.Atoi(match[3])
	return uint64(kind*1e9 + id*1e1 + ver), nil
}

// Parts extracts the kind number, identifier number and verifier digit
// parts out of the provided CUIT number and returns them as the three
// return values.
//
// This functions discards any decimal digit exceeding the allowed range.
func Parts(cuit uint64) (kind uint8, id uint32, ver uint8) {
	clean := cuit % 1e11 // get rid of possible excess digits
	return uint8(clean / 1e9), uint32((clean % 1e9) / 1e1), uint8(clean % 1e1)
}

// Format returns a standard formatted string representation of the provided
// CUIT number.
//
// This function uses Parts(cuit) to extract the constituent parts of the
// CUIT number, hence that function behavior regarding digits exceeding
// the allowed range is maintained.
func Format(cuit uint64) string {
	kind, id, ver := Parts(cuit)
	return fmt.Sprintf("%02d-%08d-%01d", kind, id, ver)
}
