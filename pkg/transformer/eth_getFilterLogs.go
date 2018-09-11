package transformer

// ProxyETHGetFilterLogs implements ETHProxy
type ProxyETHGetFilterLogs struct {
	*ProxyETHGetFilterChanges
}

func (p *ProxyETHGetFilterLogs) Method() string {
	return "eth_getFilterLogs"
}
