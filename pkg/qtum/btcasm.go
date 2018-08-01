package qtum

import (
	"math/big"
	"strings"

	"github.com/pkg/errors"
)

type (
	// ASM is Bitcoin Script extended by Qtum to support smart contracts
	ASM struct {
		VMVersion   string
		GasLimitStr string
		GasPriceStr string
		Instructor  string
	}
	CallASM struct {
		ASM
		callData        string
		ContractAddress string
	}
	CreateASM struct {
		ASM
		callData string
	}
)

func (asm *ASM) GasPrice() (*big.Int, error) {
	return stringNumberToBigInt(asm.GasPriceStr)
}

func (asm *ASM) GasLimit() (*big.Int, error) {
	return stringNumberToBigInt(asm.GasLimitStr)
}

func (asm *CreateASM) CallData() string {
	return asm.callData
}

func (asm *CallASM) CallData() string {
	return asm.callData
}

func ParseCreateASM(asm string) (*CreateASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 5 {
		return nil, errors.New("invalid create ASM")
	}

	return &CreateASM{
		ASM: ASM{
			VMVersion:   parts[0],
			GasLimitStr: parts[1],
			GasPriceStr: parts[2],
			Instructor:  parts[4],
		},
		callData: parts[3],
	}, nil
}

func ParseCallASM(asm string) (*CallASM, error) {
	parts := strings.Split(asm, " ")
	if len(parts) < 6 {
		return nil, errors.New("invalid call ASM")
	}

	return &CallASM{
		ASM: ASM{
			VMVersion:   parts[0],
			GasLimitStr: parts[1],
			GasPriceStr: parts[2],
			Instructor:  parts[5],
		},
		callData:        parts[3],
		ContractAddress: parts[4],
	}, nil
}

func stringNumberToBigInt(str string) (*big.Int, error) {
	var success bool
	v := new(big.Int)
	if v, success = v.SetString(str, 10); !success {
		return nil, errors.Errorf("Failed to parse big.Int: %s", str)
	}
	return v, nil
}
