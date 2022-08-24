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
	teamChannel   chan model.Team
	playerChannel chan model.Player
	mu            sync.Mutex
	done          chan bool
}

func NewRendererService(teamChannel chan model.Team, playerChannel chan model.Player, done chan bool, config *cfg.Config) (model.Runner, error) {

	return &renderer{
		config:        config,
		teamChannel:   teamChannel,
		done:          done,
		playerChannel: playerChannel,
	}, nil

}

func (r *renderer) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	var players = make(map[string]model.Player)
	for i := 0; i < r.config.MaxRenderWorkers; i++ {
		wg.Add(1)
		go r.worker(&wg, players)
	}
	wg.Wait()
	for _, p := range players {
		fmt.Printf("%v \n", p)
	}
	r.done <- true
	return nil

}

func (r *renderer) worker(wg *sync.WaitGroup, players map[string]model.Player) {
	defer wg.Done()
	for team := range r.teamChannel {
		for _, player := range team.Data.Team.Players {
			var p model.Player
			if players[player.Id].Id == "" {
				if team.Data.Team.IsNational {
					p.Id = player.Id
					p.Name = player.Name
					p.Country = player.Country
					p.Age = player.Age
					r.mu.Lock()
					players[player.Id] = p
					r.mu.Unlock()
				} else {
					p.Id = player.Id
					p.Name = player.Name
					p.Team = team.Data.Team.Name
					p.Age = player.Age
					r.mu.Lock()
					players[p.Id] = p
					r.mu.Unlock()
				}
			} else {
				if team.Data.Team.IsNational {
					r.mu.Lock()
					a := players[player.Id]
					a.Country = player.Country
					players[player.Id] = a
					r.mu.Unlock()
				} else {
					r.mu.Lock()
					a := players[player.Id]
					a.Team = team.Data.Team.Name
					players[player.Id] = a
					r.mu.Unlock()
				}
			}

		}
	}
}
