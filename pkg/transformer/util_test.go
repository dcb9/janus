package transformer

import (
	"testing"
)

func TestEthValueToQtumAmount(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"in":   "0x64",
			"want": float64(0.000001),
		},
		{

			"in":   "0x1",
			"want": 0.00000001,
		},
	}
	for _, c := range cases {
		in := c["in"].(string)
		want := c["want"].(float64)
		got, err := EthValueToQtumAmount(in)
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf("in: %s, want: %f, got: %f", in, want, got)
		}
	}
}

func TestQtumAmountToEthValue(t *testing.T) {
	in, want := 0.000001, "0x64"
	got, err := QtumAmountToEthValue(in)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("in: %f, want: %s, got: %s", in, want, got)
	}
}
