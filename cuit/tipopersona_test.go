// MIT License
//
// Copyright (c) 2019 Pablo Ignacio Lalloni
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

import "fmt"

func ExampleTipoPersonaCUIT() {
	// using switch/case
	show := func(c uint64) {
		switch TipoPersonaCUIT(c) {
		case PersonaFísica:
			fmt.Println("física")
			// or do anything...
		case PersonaJurídica:
			fmt.Println("jurídica")
			// or do anything...
		}
	}
	show(30123456781)
	show(20123456781)
	// Stringer implementation
	fmt.Println(TipoPersonaCUIT(30123456781))
	fmt.Println(TipoPersonaCUIT(20123456781))
	// Output:
	// jurídica
	// física
	// Persona Jurídica
	// Persona Física
}
