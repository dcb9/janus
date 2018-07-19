package eth

// eth_sendTransaction
type TransactionReq struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`      // optional
	GasPrice string `json:"gasPrice"` // optional
	Value    string `json:"value"`    // optional
	Data     string `json:"data"`     // optional
	Nonce    string `json:"nonce"`    // optional
}

// see: https://ethereum.stackexchange.com/questions/8384/transfer-an-amount-between-two-ethereum-accounts-using-json-rpc
func (t *TransactionReq) IsSendEther() bool {
	// data must be empty
	return t.Value != "" && t.To != "" && t.From != "" && t.Data == ""
}

func (t *TransactionReq) IsCreateContract() bool {
	return t.To == "" && t.Data != ""
}

func (t *TransactionReq) IsCallContract() bool {
	return t.To != "" && t.Data != ""
}

func (t *TransactionReq) GetGas() string {
	return t.Gas
}

func (t *TransactionReq) GetGasPrice() string {
	return t.GasPrice
}

// eth_call
type TransactionCallReq struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Gas      string `json:"gas"`      // optional
	GasPrice string `json:"gasPrice"` // optional
	Value    string `json:"value"`    // optional
	Data     string `json:"data"`     // optional
}

func (t *TransactionCallReq) GetGas() string {
	return t.Gas
}

func (t *TransactionCallReq) GetGasPrice() string {
	return t.GasPrice
}
