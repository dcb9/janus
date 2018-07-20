package qtum

const (
	Version = "1.0"
)

const (
	MethodGethexaddress         = "gethexaddress"
	MethodFromhexaddress        = "fromhexaddress"
	MethodSendtocontract        = "sendtocontract"
	MethodGettransactionreceipt = "gettransactionreceipt"
	MethodGettransaction        = "gettransaction"
	MethodCreatecontract        = "createcontract"
	MethodSendtoaddress         = "sendtoaddress"
	MethodCallcontract          = "callcontract"
)

type Log struct {
	Address string   `json:"address"`
	Topics  []string `json:"topics"`
	Data    string   `json:"data"`
}
