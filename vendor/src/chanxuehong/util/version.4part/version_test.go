package version

import "testing"

func TestNew(t *testing.T) {
	have := New(1, 2, 3, 4)
	want := Version{1, 2, 3, 4}
	if have != want {
		t.Errorf("Test New(1,2,3,4) failed, have %+v, want %+v", have, want)
		return
	}
}

func TestString(t *testing.T) {
	v := Version{1, 2, 3, 4}
	have := v.String()
	want := "1.2.3.4"
	if have != want {
		t.Errorf("Version{1, 2, 3, 4}.String() failed, have %q, want %q", have, want)
		return
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		str string
		v   Version
		ok  bool
	}{
		{
			"1",
			Version{1, 0, 0, 0},
			true,
		},
		{
			"1.",
			Version{1, 0, 0, 0},
			true,
		},
		{
			"1.2",
			Version{1, 2, 0, 0},
			true,
		},
		{
			"1.2.",
			Version{1, 2, 0, 0},
			true,
		},
		{
			"1.2.3",
			Version{1, 2, 3, 0},
			true,
		},
		{
			"1.2.3.",
			Version{1, 2, 3, 0},
			true,
		},
		{
			"1.2.3.4",
			Version{1, 2, 3, 4},
			true,
		},

		{
			"1.2.3.4.",
			Version{1, 2, 3, 4},
			false,
		},
		{
			"",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.3.4.5",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.3.4.5.6",
			Version{0, 0, 0, 0},
			false,
		},

		{
			".1",
			Version{0, 0, 0, 0},
			false,
		},
		{
			".1.2",
			Version{0, 0, 0, 0},
			false,
		},
		{
			".1.2.3",
			Version{0, 0, 0, 0},
			false,
		},
		{
			".1.2.3.4",
			Version{0, 0, 0, 0},
			false,
		},

		{
			"a",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.a",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.a",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.3.a",
			Version{0, 0, 0, 0},
			false,
		},

		{
			"1..",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2..",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.3..",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.3.4..",
			Version{0, 0, 0, 0},
			false,
		},

		{
			"1..2",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2..3",
			Version{0, 0, 0, 0},
			false,
		},
		{
			"1.2.3..4",
			Version{0, 0, 0, 0},
			false,
		},
	}

	for _, item := range tests {
		v, ok := Parse(item.str)
		if ok != item.ok {
			t.Errorf("Parse(%q) failed, have(Version, %v), want(Version, %v)", item.str, ok, item.ok)
			return
		}
		if ok && v != item.v {
			t.Errorf("Parse(%q) failed, have(%+v, %v), want(%+v, %v)", item.str, v, ok, item.v, item.ok)
			return
		}
	}
}
