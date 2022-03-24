package configuration

var RuntimeConf = RuntimeConfig{}

type RuntimeConfig struct {
	Discord Discord `yaml:"discord"`
	Api     Api     `yaml:"api"`
}

type Discord struct {
	BotToken string `yaml:"botToken"`
}

type Api struct {
	OsmosisUrl     string `yaml:"omosisUrl"`
	OsmosisApiKey  string `yaml:"omosisApiKey"`
	OsmosisAddress string `yaml:"osmosisAddress"`
}
