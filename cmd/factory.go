package cmd

import (
	"fmt"
	"log"
	cfg "project/config"
	"project/pkg/model"
	"project/pkg/service"
)

type factory struct {
	PipelineStages []model.Runner
	Response       chan model.PlayerData
}

func NewFactory() (*factory, error) {
	var idChannel = make(chan int, 1000)
	var playerChannel = make(chan model.Player, 500)
	var responseChannel = make(chan model.PlayerData, 500)
	config := cfg.NewConfig()

	generatorService, err1 := service.NewGeneratorService(idChannel, config)
	if err1 != nil {
		log.Fatal(fmt.Sprintf("establish generate service : %d", err1))
	}

	dataGetterService, err2 := service.NewDataGetterService(idChannel, playerChannel, config)
	if err2 != nil {
		log.Fatal(fmt.Sprintf("establish getter service : %d", err2))
	}

	rendererService, err3 := service.NewRendererService(playerChannel, responseChannel, config)
	if err3 != nil {
		log.Fatal(fmt.Sprintf("establish render service : %d", err3))
	}
	return &factory{
		PipelineStages: []model.Runner{generatorService, dataGetterService, rendererService},
		Response:       responseChannel,
	}, nil
}
