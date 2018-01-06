package simplexer

import "testing"

func TestPositionString(t *testing.T) {
	if s := (Position{ Line: 0, Column: 1 }).String(); s != "[line:0, column:1]" {
		t.Errorf("failed convert to string: excepted [line:0, column:1] but got %#v", s)
	}

	if s := (Position{ Line: 5, Column: 3 }).String(); s != "[line:5, column:3]" {
		t.Errorf("failed convert to string: excepted [line:5, column:3] but got %#v", s)
	}
}

func TestPositionCompare(t *testing.T) {
	a := Position{ Line: 0, Column: 0}
	a2 := Position{ Line: 0, Column: 0}
	b := Position{ Line: 0, Column: 5}
	c := Position{ Line: 1, Column: 3}

	if a != a2 {
		t.Errorf("Position reports %s != %s", a, a2)
	}

	if !a.Before(b) {
		t.Errorf("Position reports %s is not before of %s", a, b)
	}

	if !a.Before(c) {
		t.Errorf("Position reports %s is not before of %s", a, c)
	}

	if !b.Before(c) {
		t.Errorf("Position reports %s is not before of %s", b, c)
	}

	if !b.After(a) {
		t.Errorf("Position reports %s is not after of %s", b, a)
	}

	if !c.After(a) {
		t.Errorf("Position reports %s is not after of %s", c, a)
	}

	if !c.After(b) {
		t.Errorf("Position reports %s is not after of %s", c, b)
	}
}
