package configuration

type Configuration struct {
	App     *AppConfig   `yaml:"app"`
	JobSite *JobSite     `yaml:"jobSite"`
	Mongo   *MongoConfig `yaml:"mongo"`
	Geo     *GeoConfig   `yaml:"geo"`
}

type JobSite struct {
	BaseUrl          string `yaml:"base_url"`
	FullCardSelector string `yaml:"full_card_selector"`
	CardInfoSelector string `yaml:"card_info_selector"`
}

type MongoConfig struct {
	Hostname            string `yaml:"hostname"`
	Database            string `yaml:"database"`
	JobPostsCollection  string `yaml:"job_posts_collection"`
	GeocodingCollection string `yaml:"geocoding_collection"`
	Username            string `yaml:"username"`
	Password            string `yaml:"password"`
	RetryWrites         bool   `yaml:"retry_writes"`
}

type AppConfig struct {
	DefaultDate string `yaml:"default_date"`
	ScanFactor  int    `yaml:"scan_factor"`
}

type GeoConfig struct {
	AccessKey string `yaml:"access_key"`
	BaseUrl   string `yaml:"base_url"`
	Timeout   int    `yaml:"timeout"`
}
