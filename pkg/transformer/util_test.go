package transformer

import (
	"math/big"
	"testing"

	"github.com/dcb9/janus/pkg/eth"
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

func TestEthGasToQtum(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"in": &eth.TransactionReq{
				Gas:      "0x64",
				GasPrice: "0x1",
			},
			"wantGas":      big.NewInt(100),
			"wantGasPrice": "0.00000001",
		},
		{
			"in": &eth.TransactionReq{
				Gas:      "0x1",
				GasPrice: "0xff",
			},
			"wantGas":      big.NewInt(1),
			"wantGasPrice": "0.00000255",
		},
		{
			"in": &eth.TransactionReq{
				Gas:      "0x1",
				GasPrice: "0x64",
			},
			"wantGas":      big.NewInt(1),
			"wantGasPrice": "0.00000100",
		},
	}

	for _, c := range cases {
		in := c["in"].(*eth.TransactionReq)
		wantGas, wantGasPrice := c["wantGas"].(*big.Int), c["wantGasPrice"].(string)
		gotGas, gotGasPrice, err := EthGasToQtum(in)
		if err != nil {
			t.Error(err)
		}
		if gotGas.Cmp(wantGas) != 0 {
			t.Errorf("get Gas error in: %s, want: %d, got: %d", in.Gas, wantGas, gotGas.Int64())
		}
		if gotGasPrice != wantGasPrice {
			t.Errorf("get GasPrice error in: %s, want: %s, got: %s", in.GasPrice, wantGasPrice, gotGasPrice)
		}
	}
}
