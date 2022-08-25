package service

import (
	"context"
	"fmt"
	cfg "project/config"
	"project/pkg/model"
	"sync"
)

type renderer struct {
	config        *cfg.Config
	playerChannel chan model.Player
	players       sync.Map
	done          chan bool
}

func NewRendererService(playerChannel chan model.Player, done chan bool, config *cfg.Config) (model.Runner, error) {

	return &renderer{
		config:        config,
		done:          done,
		playerChannel: playerChannel,
		players:       sync.Map{},
	}, nil

}

func (r *renderer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	for i := 0; i < r.config.MaxRenderWorkers; i++ {
		wg.Add(1)
		go r.worker(&wg)
	}
	wg.Wait()
	r.players.Range(func(key, value interface{}) bool {
		player := value.(model.PlayerData)
		fmt.Printf("%v \n", player)
		return true
	})
	r.done <- true
	return nil

}

func (r *renderer) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for player := range r.playerChannel {
		if value, ok := r.players.Load(player.Id); ok {
			playerData := value.(model.PlayerData)
			playerData.Teams = append(playerData.Teams, player.Team)
			r.players.Store(player.Id, playerData)
		} else {
			r.players.Store(player.Id, model.PlayerData{
				Id:    player.Id,
				Name:  player.Name,
				Age:   player.Age,
				Teams: []string{player.Team},
			})
		}
	}
}
