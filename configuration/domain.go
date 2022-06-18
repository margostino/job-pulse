package configuration

type Configuration struct {
	App     *AppConfig   `yaml:"app"`
	JobSite *JobSite     `yaml:"jobSite"`
	Mongo   *MongoConfig `yaml:"mongo"`
}

type JobSite struct {
	BaseUrl          string `yaml:"baseUrl"`
	FullCardSelector string `yaml:"fullCardSelector"`
	CardInfoSelector string `yaml:"cardInfoSelector"`
}

type MongoConfig struct {
	Hostname    string `yaml:"hostname"`
	Database    string `yaml:"database"`
	Collection  string `yaml:"collection"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	RetryWrites bool   `yaml:"retryWrites"`
}

type AppConfig struct {
	DefaultDate string `yaml:"defaultDate"`
	ScanFactor  int    `yaml:"scanFactor"`
}
