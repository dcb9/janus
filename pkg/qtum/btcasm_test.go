package qtum

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseCallASM(t *testing.T) {
	samStr := "4 250000 40 60fe47b10000000000000000000000000000000000000000000000000000000000000002 cd20af1f2d6ac4173f9464030e7cef40bf9cb7c4 OP_CALL"
	got, err := ParseCallASM(samStr)
	if err != nil {
		t.Error(err)
	}
	want := &CallASM{
		ASM: ASM{
			VMVersion:   "4",
			GasLimitStr: "250000",
			GasPriceStr: "40",
			Instructor:  "OP_CALL",
		},
		callData:        "60fe47b10000000000000000000000000000000000000000000000000000000000000002",
		ContractAddress: "cd20af1f2d6ac4173f9464030e7cef40bf9cb7c4",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"parse transaction call sam error\ninput: %s\nwant: %s\ngot: %s",
			samStr,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}

func TestParseCreateASM(t *testing.T) {
	samStr := "4 6721975 40 608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b8100290000000000000000000000000000000000000000000000000000000000000001 OP_CREATE"
	got, err := ParseCreateASM(samStr)
	if err != nil {
		t.Error(err)
	}
	want := &CreateASM{
		ASM: ASM{
			VMVersion:   "4",
			GasLimitStr: "6721975",
			GasPriceStr: "40",
			Instructor:  "OP_CREATE",
		},
		callData: "608060405234801561001057600080fd5b506040516020806100f2833981016040525160005560bf806100336000396000f30060806040526004361060485763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166360fe47b18114604d5780636d4ce63c146064575b600080fd5b348015605857600080fd5b5060626004356088565b005b348015606f57600080fd5b506076608d565b60408051918252519081900360200190f35b600055565b600054905600a165627a7a7230582049a087087e1fc6da0b68ca259d45a2e369efcbb50e93f9b7fa3e198de6402b8100290000000000000000000000000000000000000000000000000000000000000001",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf(
			"parse transaction create sam error\ninput: %s\nwant: %s\ngot: %s",
			samStr,
			string(mustMarshalIndent(want, "", "  ")),
			string(mustMarshalIndent(got, "", "  ")),
		)
	}
}

func mustMarshalIndent(v interface{}, prefix, indent string) []byte {
	res, err := json.MarshalIndent(v, prefix, indent)
	if err != nil {
		panic(err)
	}
	return res
}
