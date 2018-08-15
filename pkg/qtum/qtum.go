package qtum

import (
	"github.com/dcb9/janus/pkg/qtum/insight"
	"github.com/dcb9/janus/pkg/utils"
	"github.com/pkg/errors"
)

type Qtum struct {
	*Client
	*Method
	chain   string
	Insight *insight.Insight
}

const (
	ChainMain    = "main"
	ChainTest    = "test"
	ChainRegTest = "regtest"
)

var AllChains = []string{ChainMain, ChainRegTest, ChainTest}

func New(c *Client, chain string, insight *insight.Insight) (*Qtum, error) {
	if !utils.InStrSlice(AllChains, chain) {
		return nil, errors.New("invalid qtum chain")
	}

	return &Qtum{
		Client:  c,
		Method:  &Method{Client: c},
		chain:   chain,
		Insight: insight,
	}, nil
}

func (c *Qtum) Chain() string {
	return c.chain
}
