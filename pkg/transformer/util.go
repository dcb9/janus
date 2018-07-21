package transformer

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type EthGas interface {
	GetGas() string
	GetGasPrice() string
}

func EthGasToQtum(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = big.NewInt(2500000)
	if gas := g.GetGas(); gas != "" {
		gasLimit, err = hexutil.DecodeBig(AddHexPrefix(gas))
		if err != nil {
			err = errors.Wrap(err, "decode gas")
			return
		}
	}

	gasPriceFloat64, err := EthValueToQtumAmount(g.GetGasPrice())
	if err != nil {
		return nil, "0.0", err
	}
	gasPrice = fmt.Sprintf("%f", gasPriceFloat64)

	return
}

func EthValueToQtumAmount(val string) (float64, error) {
	ethVal, err := hexutil.DecodeBig(AddHexPrefix(val))
	if err != nil {
		return 0.0, err
	}

	ethValFloat64 := new(big.Float)
	ethValFloat64, success := ethValFloat64.SetString(ethVal.String())
	if !success {
		return 0.0, errors.New("big.Float#SetString is not success")
	}

	amount := ethValFloat64.Mul(ethValFloat64, big.NewFloat(float64(1e-8)))
	result, _ := amount.Float64()

	return result, nil
}

func QtumAmountToEthValue(amount float64) (string, error) {
	bigAmount := big.NewFloat(amount)
	bigAmount = bigAmount.Mul(bigAmount, big.NewFloat(float64(1e8)))

	result := new(big.Int)
	result, success := result.SetString(bigAmount.String(), 10)
	if !success {
		return "0x0", errors.New("big.Int#SetString is not success")
	}

	return hexutil.EncodeBig(result), nil
}

func RemoveHexPrefix(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex[2:]
	}
	return hex
}

func IsEthHexAddress(str string) bool {
	return strings.HasPrefix(str, "0x") || common.IsHexAddress("0x"+str)
}

func AddHexPrefix(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex
	}
	return "0x" + hex
}

func unmarshalRequest(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return UnmarshalRequestErr
	}
	return nil
}
