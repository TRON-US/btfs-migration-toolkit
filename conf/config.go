package conf

type Config struct {
	IpfsUrl      string
	SoterUrl     string
	UploaderPath   string
	VerifierPath string
	BatchSize    int
	Logger       LogConfig
}

type LogConfig struct {
	Path  string
	Level string
}