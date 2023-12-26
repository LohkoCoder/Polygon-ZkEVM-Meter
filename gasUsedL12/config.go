package gasUsedL12

var layer1Network = struct {
	Name           string
	URL            string
	ChainID        uint64
	PolygonPoEAddr string
}{
	"L1 Mainnet",
	"wss://mainnet.infura.io/ws/v3/a2816bc61057481188c0a772cee3134e",
	1,
	"0x5132A183E9F3CB7C848b0AAC5Ae0c4f0491B7aB2",
}

var layer2Network = struct {
	Name    string
	URL     string
	ChainID uint64
}{
	"PolygonZkEVM Mainnet",
	"https://zkevm-rpc.com",
	1101,
}
