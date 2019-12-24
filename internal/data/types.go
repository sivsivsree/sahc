package data

type Configuration struct {
	Version  float64 `yaml:"version"`
	Services []struct {
		Name     string `yaml:"name"`
		Interval int    `yaml:"interval"`
	} `yaml:"services"`
}


