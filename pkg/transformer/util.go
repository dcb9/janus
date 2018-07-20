package transformer

import (
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
		gasLimit, err = hexutil.DecodeBig(gas)
		if err != nil {
			err = errors.Wrap(err, "decode gas")
			return
		}
	}
	gasPrice = "0.0000004"
	// fixme parse gas price

	return
}

func EthHexToQtum(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex[2:]
	}
	return hex
}

func IsEthHex(str string) bool {
	return strings.HasPrefix(str, "0x") || common.IsHexAddress("0x"+str)
}

func QtumHexToEth(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex
	}
	return "0x" + hex
}
