package transformer

import "testing"

func TestEthValueToQtumAmount(t *testing.T) {
	in := "0x64"
	want := 0.000001
	got, err := EthValueToQtumAmount(in)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("in: %s, want: %f, got: %f", in, want, got)
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
