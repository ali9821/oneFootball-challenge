package cfg

import (
	"github.com/spf13/viper"
	"math"
)

type Config struct {
	OutputFormat     string
	OutputFile       string
	OutputPath       string
	MaxWorkers       int
	MaxRenderWorkers int
	MaxFetch         int
	Api              string
	WantedTeams      []string
}

func NewConfig() *Config {

	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.SetConfigType("yml")

	setInitialConfig(v)

	return &Config{
		OutputFormat:     v.GetString("output_format"),
		OutputFile:       v.GetString("output_file"),
		OutputPath:       v.GetString("output_path"),
		MaxWorkers:       v.GetInt("max_getter_worker_size"),
		MaxRenderWorkers: v.GetInt("max_renderer_worker_size"),
		MaxFetch:         v.GetInt("max_fetch_size"),
		Api:              v.GetString("api_path"),
		WantedTeams:      v.GetStringSlice("wanted_teams"),
	}

}

func setInitialConfig(v *viper.Viper) {
	v.SetDefault("output_format", "stdout")
	v.SetDefault("output_file", "result.txt")
	v.SetDefault("output_path", ".")
	v.SetDefault("max_getter_worker_size", 100)
	v.SetDefault("max_renderer_worker_size", 100)
	v.SetDefault("max_fetch_size", math.MaxInt)
	v.SetDefault("api_path", "https://api-origin.onefootball.com/score-one-proxy/api/teams/en/%d.json")
	v.SetDefault("wanted_teams", []string{"Germany", "England", "France", "Spain", "Manchester United", "Arsenal", "Chelsea", "Barcelona", "Real Madrid", "Bayern Munich"})
}
