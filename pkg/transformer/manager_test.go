package transformer

import "testing"

func TestRemoveHexPrefix(t *testing.T) {
	cases := []map[string]string{
		{
			"in":   "0x8d124864e8840a114a8772c1daf82b61eb4dca01",
			"want": "8d124864e8840a114a8772c1daf82b61eb4dca01",
		},
		{
			"in":   "8d124864e8840a114a8772c1daf82b61eb4dca01",
			"want": "8d124864e8840a114a8772c1daf82b61eb4dca01",
		},
		{
			"in":   "",
			"want": "",
		},
	}

	for _, c := range cases {
		in, want := c["in"], c["want"]
		if got := RemoveHexPrefix(in); got != want {
			t.Errorf("err: in: %s, want: %s, got: %s", in, want, got)
		}
	}
}
