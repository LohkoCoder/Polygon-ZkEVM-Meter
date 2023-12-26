package utils

var Networks = []struct {
	Name       string
	URL        string
	ChainID    uint64
	PrivateKey string
}{
	{Name: "Local L2", URL: "http://106.75.76.70:8123", ChainID: 10898, PrivateKey: "0x4689d34373cfd5d0b3cc8602352e3c45cb5409e57d0017ada3597ad4f8262155"},
	{Name: "Local L1", URL: "http://106.75.76.70:8545", ChainID: 10898, PrivateKey: "0x1010cfa44a9ada50594a833a997c8b6e61461e6150b20964d2f5b497e7f47e56"},
	{Name: "Local L1", URL: "http://106.75.76.70:8545", ChainID: 10898, PrivateKey: "0xe817a8de9d7e7d55df8b7e0c9add6a64d27f85c8423401707efa23037788c2fc"},
}
