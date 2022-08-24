package service

import (
	"context"
	cfg "project/config"
	"project/pkg/model"
)

type generator struct {
	idChannel chan int
	config    *cfg.Config
}

func NewGeneratorService(idChannel chan int, config *cfg.Config) (model.Runner, error) {
	return &generator{
		idChannel: idChannel,
		config:    config,
	}, nil

}

func (g *generator) Run(ctx context.Context) error {
	defer close(g.idChannel)
	for i := 1; i < g.config.MaxFetch; i++ {
		g.idChannel <- i
	}
	return nil
}
