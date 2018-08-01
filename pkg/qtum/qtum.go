package qtum

type Qtum struct {
	*Client
	*Method
}

func New(c *Client) *Qtum {
	return &Qtum{
		Client: c,
		Method: &Method{Client: c},
	}
}
