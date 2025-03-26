package config

type Configuration struct {
	AllowedOrigins           []string `split_words:"true" default:"(.)+.localhost:([0-9]+)?"`
	HttpProxy                string   `split_words:"true"`
	HttpsProxy               string   `split_words:"true"`
	NoProxy                  string   `split_words:"true"`
	HTTPPort                 string   `envconfig:"HTTP_PORT" default:"9090"`
	DbServer                 string   `required:"true" split_words:"true" default:"mongodb://localhost:27017"`
	DbUser                   string   `required:"true" split_words:"true" default:"admin"`
	DbPassword               string   `required:"true" split_words:"true" default:"secretpassword"`
	BppId                    string   `required:"true" split_words:"true" default:"bpp"`
	BppUri                   string   `required:"true" split_words:"true" default:"http://localhost:8080"`
	BapId                    string   `required:"true" split_words:"true" default:"bap"`
	BapUri                   string   `required:"true" split_words:"true" default:"http://localhost:9090"`
	RecommendationServiceURL string   `required:"false" split_words:"true"`
	RedisHost                string   `required:"true" split_words:"true" default:"localhost"`
	RedisPort                string   `required:"true" split_words:"true" default:"6379"`
	RedisPassword            string   `required:"true" split_words:"true" default:"yourpassword"`
}

var Config Configuration
