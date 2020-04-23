package version

import "testing"

func TestCompare(t *testing.T) {
	tests := []struct {
		a      Version
		b      Version
		result int // -1 28个, +1 28个, 0 8个, 总的 64 个
	}{
		{
			Version{0, 0, 0},
			Version{0, 0, 0},
			0,
		},
		{
			Version{0, 0, 0},
			Version{0, 0, 1},
			-1,
		},
		{
			Version{0, 0, 0},
			Version{0, 1, 0},
			-1,
		},
		{
			Version{0, 0, 0},
			Version{0, 1, 1},
			-1,
		},
		{
			Version{0, 0, 0},
			Version{1, 0, 0},
			-1,
		},
		{
			Version{0, 0, 0},
			Version{1, 0, 1},
			-1,
		},
		{
			Version{0, 0, 0},
			Version{1, 1, 0},
			-1,
		},
		{
			Version{0, 0, 0},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 0},
			1,
		},
		{
			Version{0, 0, 1},
			Version{0, 0, 1},
			0,
		},
		{
			Version{0, 0, 1},
			Version{0, 1, 0},
			-1,
		},
		{
			Version{0, 0, 1},
			Version{0, 1, 1},
			-1,
		},
		{
			Version{0, 0, 1},
			Version{1, 0, 0},
			-1,
		},
		{
			Version{0, 0, 1},
			Version{1, 0, 1},
			-1,
		},
		{
			Version{0, 0, 1},
			Version{1, 1, 0},
			-1,
		},
		{
			Version{0, 0, 1},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{0, 1, 0},
			Version{0, 0, 0},
			1,
		},
		{
			Version{0, 1, 0},
			Version{0, 0, 1},
			1,
		},
		{
			Version{0, 1, 0},
			Version{0, 1, 0},
			0,
		},
		{
			Version{0, 1, 0},
			Version{0, 1, 1},
			-1,
		},
		{
			Version{0, 1, 0},
			Version{1, 0, 0},
			-1,
		},
		{
			Version{0, 1, 0},
			Version{1, 0, 1},
			-1,
		},
		{
			Version{0, 1, 0},
			Version{1, 1, 0},
			-1,
		},
		{
			Version{0, 1, 0},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{0, 1, 1},
			Version{0, 0, 0},
			1,
		},
		{
			Version{0, 1, 1},
			Version{0, 0, 1},
			1,
		},
		{
			Version{0, 1, 1},
			Version{0, 1, 0},
			1,
		},
		{
			Version{0, 1, 1},
			Version{0, 1, 1},
			0,
		},
		{
			Version{0, 1, 1},
			Version{1, 0, 0},
			-1,
		},
		{
			Version{0, 1, 1},
			Version{1, 0, 1},
			-1,
		},
		{
			Version{0, 1, 1},
			Version{1, 1, 0},
			-1,
		},
		{
			Version{0, 1, 1},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{1, 0, 0},
			Version{0, 0, 0},
			1,
		},
		{
			Version{1, 0, 0},
			Version{0, 0, 1},
			1,
		},
		{
			Version{1, 0, 0},
			Version{0, 1, 0},
			1,
		},
		{
			Version{1, 0, 0},
			Version{0, 1, 1},
			1,
		},
		{
			Version{1, 0, 0},
			Version{1, 0, 0},
			0,
		},
		{
			Version{1, 0, 0},
			Version{1, 0, 1},
			-1,
		},
		{
			Version{1, 0, 0},
			Version{1, 1, 0},
			-1,
		},
		{
			Version{1, 0, 0},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{1, 0, 1},
			Version{0, 0, 0},
			1,
		},
		{
			Version{1, 0, 1},
			Version{0, 0, 1},
			1,
		},
		{
			Version{1, 0, 1},
			Version{0, 1, 0},
			1,
		},
		{
			Version{1, 0, 1},
			Version{0, 1, 1},
			1,
		},
		{
			Version{1, 0, 1},
			Version{1, 0, 0},
			1,
		},
		{
			Version{1, 0, 1},
			Version{1, 0, 1},
			0,
		},
		{
			Version{1, 0, 1},
			Version{1, 1, 0},
			-1,
		},
		{
			Version{1, 0, 1},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{1, 1, 0},
			Version{0, 0, 0},
			1,
		},
		{
			Version{1, 1, 0},
			Version{0, 0, 1},
			1,
		},
		{
			Version{1, 1, 0},
			Version{0, 1, 0},
			1,
		},
		{
			Version{1, 1, 0},
			Version{0, 1, 1},
			1,
		},
		{
			Version{1, 1, 0},
			Version{1, 0, 0},
			1,
		},
		{
			Version{1, 1, 0},
			Version{1, 0, 1},
			1,
		},
		{
			Version{1, 1, 0},
			Version{1, 1, 0},
			0,
		},
		{
			Version{1, 1, 0},
			Version{1, 1, 1},
			-1,
		},
		{
			Version{1, 1, 1},
			Version{0, 0, 0},
			1,
		},
		{
			Version{1, 1, 1},
			Version{0, 0, 1},
			1,
		},
		{
			Version{1, 1, 1},
			Version{0, 1, 0},
			1,
		},
		{
			Version{1, 1, 1},
			Version{0, 1, 1},
			1,
		},
		{
			Version{1, 1, 1},
			Version{1, 0, 0},
			1,
		},
		{
			Version{1, 1, 1},
			Version{1, 0, 1},
			1,
		},
		{
			Version{1, 1, 1},
			Version{1, 1, 0},
			1,
		},
		{
			Version{1, 1, 1},
			Version{1, 1, 1},
			0,
		},
	}

	for _, v := range tests {
		result := Compare(v.a, v.b)
		if result != v.result {
			t.Errorf("Compare(%+v, %+v) failed, have %d, want %d", v.a, v.b, result, v.result)
			return
		}
	}
}
