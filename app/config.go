package app

type Config struct {
	UploadUrl string
	TmpFolder string
	AccessKey string
}

var (
	config *Config
)

func init() {
	config = new(Config)
}

func GetConfig() *Config {
	return config
}
