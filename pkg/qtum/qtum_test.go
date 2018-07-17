package qtum

import "testing"

func TestGetHexAddress(t *testing.T) {
}

func TestFromHexAddress(t *testing.T) {
	cases := []map[string]string{
		//{
		//	"in":   "dd2c6512563e4274dafd8312e0e738ede48f3046",
		//	"want": "qdiqg2mp646KhSQjVud3whv6C34hNHQnL2",
		//},
	}
	for _, c := range cases {
		in, want := c["in"], c["want"]

		got, err := FromHexAddress(in)
		if err != nil {
			t.Error("err", err)
		}

		if got != want {
			t.Errorf("gethexaddress error in: %s, want: %s, got: %s", in, want, got)
		}
	}
}
