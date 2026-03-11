package string

import "testing"

func TestConvertYuanStringToFen(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{in: "12.34", want: 1234},
		{in: "12.345", want: 1235},
		{in: "12.344", want: 1234},
		{in: "0.005", want: 1},
		{in: "-0.005", want: -1},
		{in: "100", want: 10000},
		{in: " 1.2 ", want: 120},
		{in: "", want: 0},
		{in: "abc", want: 0},
	}

	for _, c := range cases {
		got := ConvertYuanStringToFen(c.in)
		if got != c.want {
			t.Fatalf("ConvertYuanStringToFen(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestConvertFenToYuanString(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{in: 1234, want: "12.34"},
		{in: 120, want: "1.20"},
		{in: 1, want: "0.01"},
		{in: 0, want: "0.00"},
		{in: -1, want: "-0.01"},
		{in: -1234, want: "-12.34"},
	}

	for _, c := range cases {
		got := ConvertFenToYuanString(c.in)
		if got != c.want {
			t.Fatalf("ConvertFenToYuanString(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}

