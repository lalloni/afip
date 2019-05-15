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

package periodo

import "testing"

func TestParse(t *testing.T) {
	type args struct {
		t tipoPeriodo
		v string
	}
	tests := []struct {
		name      string
		args      args
		wantMatch bool
		wantY     uint
		wantM     uint
		wantD     uint
	}{
		{"base diario", args{Diario, "20000101"}, true, 2000, 1, 1},
		{"base diario ceros", args{Diario, "20000000"}, true, 2000, 0, 0},
		{"base mensual", args{Mensual, "200001"}, true, 2000, 1, 0},
		{"base anual", args{Anual, "2000"}, true, 2000, 0, 0},
		{"mal diario", args{Diario, "2000"}, false, 0, 0, 0},
		{"mal mensual", args{Mensual, "2000"}, false, 0, 0, 0},
		{"mal mensual mes", args{Mensual, "200014"}, false, 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, gotY, gotM, gotD := Parse(tt.args.t, tt.args.v)
			if gotMatch != tt.wantMatch {
				t.Errorf("Parse() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
			if tt.wantY != 0 && gotY != tt.wantY {
				t.Errorf("Parse() gotY = %v, want %v", gotY, tt.wantY)
			}
			if tt.wantM != 0 && gotM != tt.wantM {
				t.Errorf("Parse() gotM = %v, want %v", gotM, tt.wantM)
			}
			if tt.wantD != 0 && gotD != tt.wantD {
				t.Errorf("Parse() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
	t.Run("panic con mal tipo", func(t *testing.T) {
		defer func() {
			s := recover()
			if s == nil {
				t.Error("Parse() expected panic did not occur")
			} else if s != periododesconocido {
				t.Errorf("Parse() panic got %q, want %q", s, periododesconocido)
			}
		}()
		Parse(tipoPeriodo(0xff), "20000101") // should panic
	})
	t.Run("mal entero", func(t *testing.T) {
		if ok, _, _, _ := Parse(Diario, "blah"); ok {
			t.Error("Parse() expected false")
		}
	})
}

func TestComposePeriodoDiario(t *testing.T) {
	type args struct {
		y uint
		m uint
		d uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{"a", args{1000, 1, 1}, 10000101},
		{"b", args{1000, 0, 0}, 10000000},
		{"c", args{0, 0, 0}, 0},
		{"d", args{9999, 0, 0}, 99990000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposePeriodoDiario(tt.args.y, tt.args.m, tt.args.d); got != tt.want {
				t.Errorf("ComposePeriodoDiario() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecomposePeriodoDiario(t *testing.T) {
	type args struct {
		p uint
	}
	tests := []struct {
		name  string
		args  args
		wantY uint
		wantM uint
		wantD uint
	}{
		{"a", args{10000101}, 1000, 1, 1},
		{"b", args{10000000}, 1000, 0, 0},
		{"c", args{100010101}, 10001, 1, 1},
		{"d", args{9990101}, 999, 1, 1},
		{"e", args{10}, 0, 0, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotY, gotM, gotD := DecomposePeriodoDiario(tt.args.p)
			if gotY != tt.wantY {
				t.Errorf("DecomposePeriodoDiario() gotY = %v, want %v", gotY, tt.wantY)
			}
			if gotM != tt.wantM {
				t.Errorf("DecomposePeriodoDiario() gotM = %v, want %v", gotM, tt.wantM)
			}
			if gotD != tt.wantD {
				t.Errorf("DecomposePeriodoDiario() gotD = %v, want %v", gotD, tt.wantD)
			}
		})
	}
}

func TestComposePeriodoMensual(t *testing.T) {
	type args struct {
		y uint
		m uint
	}
	tests := []struct {
		name string
		args args
		want uint
	}{
		{"a", args{1000, 1}, 100001},
		{"b", args{10, 1}, 1001},
		{"c", args{9999, 40}, 999940},
		{"d", args{19999, 40}, 1999940},
		{"e", args{19999, 0}, 1999900},
		{"f", args{1000, 100}, 100100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComposePeriodoMensual(tt.args.y, tt.args.m); got != tt.want {
				t.Errorf("ComposePeriodoMensual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecomposePeriodoMensual(t *testing.T) {
	type args struct {
		p uint
	}
	tests := []struct {
		name  string
		args  args
		wantY uint
		wantM uint
	}{
		{"a", args{100001}, 1000, 1},
		{"b", args{100000}, 1000, 0},
		{"c", args{1000101}, 10001, 1},
		{"d", args{99901}, 999, 1},
		{"e", args{10}, 0, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotY, gotM := DecomposePeriodoMensual(tt.args.p)
			if gotY != tt.wantY {
				t.Errorf("DecomposePeriodoMensual() gotY = %v, want %v", gotY, tt.wantY)
			}
			if gotM != tt.wantM {
				t.Errorf("DecomposePeriodoMensual() gotM = %v, want %v", gotM, tt.wantM)
			}
		})
	}
}

func TestCheckPeriodoDiario(t *testing.T) {
	type args struct {
		y uint
		m uint
		d uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"base", args{1993, 10, 20}, true},
		{"min", args{1000, 1, 1}, true},
		{"max", args{9999, 12, 31}, true},
		{"day too big", args{2000, 13, 32}, false},
		{"day zero", args{2000, 12, 0}, true},
		{"month too big ", args{2000, 13, 1}, false},
		{"month zero day not", args{2000, 0, 1}, false},
		{"month zero day zero", args{2000, 0, 0}, true},
		{"year too big ", args{10000, 13, 1}, false},
		{"year too small ", args{999, 0, 1}, false},
		{"jan has 31 days", args{2000, 1, 31}, true},
		{"feb 2000 had 29 days", args{2000, 2, 29}, true},
		{"feb 2001 had not 29 days", args{2001, 2, 29}, false},
		{"mar has 31 days", args{2000, 3, 31}, true},
		{"apr has 30 days", args{2000, 4, 31}, false},
		{"may has 31 days", args{2000, 5, 31}, true},
		{"jun has 30 days", args{2000, 6, 31}, false},
		{"jul has 31 days", args{2000, 7, 31}, true},
		{"ago has 31 days", args{2000, 8, 31}, true},
		{"set has 30 days", args{2000, 9, 31}, false},
		{"oct has 31 days", args{2000, 10, 31}, true},
		{"nov has 30 days", args{2000, 11, 31}, false},
		{"dec has 31 days", args{2000, 12, 31}, true},
		{"zero", args{0, 0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPeriodoDiario(tt.args.y, tt.args.m, tt.args.d); got != tt.want {
				t.Errorf("CheckPeriodoDiario() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPeriodoDiarioCompound(t *testing.T) {
	type args struct {
		v uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"base", args{19931020}, true},
		{"min", args{10000101}, true},
		{"max", args{99991231}, true},
		{"day too big", args{20001332}, false},
		{"day zero", args{20001200}, true},
		{"month too big ", args{20001301}, false},
		{"month zero day not", args{20000001}, false},
		{"month zero day zero", args{20000000}, true},
		{"year too big ", args{100001301}, false},
		{"year too small ", args{9990001}, false},
		{"jan has 31 days", args{20000131}, true},
		{"feb 2000 had 29 days", args{20000229}, true},
		{"feb 2001 had not 29 days", args{20010229}, false},
		{"mar has 31 days", args{20000331}, true},
		{"apr has 30 days", args{20000431}, false},
		{"may has 31 days", args{20000531}, true},
		{"jun has 30 days", args{20000631}, false},
		{"jul has 31 days", args{20000731}, true},
		{"ago has 31 days", args{20000831}, true},
		{"set has 30 days", args{20000931}, false},
		{"oct has 31 days", args{20001031}, true},
		{"nov has 30 days", args{20001131}, false},
		{"dec has 31 days", args{20001231}, true},
		{"zero", args{0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPeriodoDiarioCompound(tt.args.v); got != tt.want {
				t.Errorf("CheckPeriodoDiarioCompound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPeriodoMensual(t *testing.T) {
	type args struct {
		y uint
		m uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"base", args{1993, 10}, true},
		{"min", args{1000, 1}, true},
		{"max", args{9999, 12}, true},
		{"too big 1", args{9999, 13}, false},
		{"too big 2", args{19999, 12}, false},
		{"too small 1", args{999, 0}, false},
		{"bad month 1", args{2000, 13}, false},
		{"bad month 2", args{2000, 50}, false},
		{"month zero", args{2000, 0}, true},
		{"zero", args{0, 0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPeriodoMensual(tt.args.y, tt.args.m); got != tt.want {
				t.Errorf("CheckPeriodoMensual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPeriodoMensualCompound(t *testing.T) {
	type args struct {
		v uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"base", args{199310}, true},
		{"min", args{100001}, true},
		{"max", args{999912}, true},
		{"too big 1", args{999913}, false},
		{"too big 2", args{1999912}, false},
		{"too small 2", args{99900}, false},
		{"bad month 1", args{200013}, false},
		{"bad month 2", args{200050}, false},
		{"month zero", args{200000}, true},
		{"zero", args{0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPeriodoMensualCompound(tt.args.v); got != tt.want {
				t.Errorf("CheckPeriodoMensualCompound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckPeriodoAnual(t *testing.T) {
	type args struct {
		y uint
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"base", args{1993}, true},
		{"min", args{1000}, true},
		{"max", args{9999}, true},
		{"too big 1", args{10000}, false},
		{"too big 2", args{19999}, false},
		{"too small 1", args{999}, false},
		{"too small 2", args{9}, false},
		{"zero", args{0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPeriodoAnual(tt.args.y); got != tt.want {
				t.Errorf("CheckPeriodoAnual() = %v, want %v", got, tt.want)
			}
		})
	}
}
