package transformer

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/dcb9/janus/pkg/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type EthGas interface {
	GasHex() string
	GasPriceHex() string
}

func EthGasToQtum(g EthGas) (gasLimit *big.Int, gasPrice string, err error) {
	gasLimit = big.NewInt(40000000)
	if gas := g.GasHex(); gas != "" {
		gasLimit, err = utils.DecodeBig(gas)
		if err != nil {
			err = errors.Wrap(err, "decode gas")
			return
		}
	}

	gasPriceFloat64, err := EthValueToQtumAmount(g.GasPriceHex())
	if err != nil {
		return nil, "0.0", err
	}
	gasPrice = fmt.Sprintf("%.8f", gasPriceFloat64)

	return
}

func EthValueToQtumAmount(val string) (float64, error) {
	if val == "" {
		return 0.0000004, nil
	}

	ethVal, err := utils.DecodeBig(val)
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

func unmarshalRequest(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return errors.Wrap(UnmarshalRequestErr, err.Error())
	}
	return nil
}
