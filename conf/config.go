package conf

type Config struct {
	IpfsUrl      string
	SoterUrl     string
	PrivateKey   string
	UserAddress  string
	BatchSize    int
	Logger       LogConfig
}

type LogConfig struct {
	Path  string
	Level string
}
