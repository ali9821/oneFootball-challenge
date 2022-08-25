package service

import (
	"context"
	"fmt"
	"net/http"
	cfg "project/config"
	"project/pkg/model"
	"project/tools"
	"sync"
)

type dataGetter struct {
	idChannel     chan int
	playerChannel chan model.Player
	wantedTeams   int
	config        *cfg.Config
	mu            sync.Mutex
}

func NewDataGetterService(idChannel chan int, playerChannel chan model.Player, config *cfg.Config) (model.Runner, error) {
	return &dataGetter{
		idChannel:     idChannel,
		config:        config,
		playerChannel: playerChannel,
		wantedTeams:   len(config.WantedTeams),
	}, nil
}

func (d *dataGetter) Run(ctx context.Context) error {
	ctx, done := context.WithCancel(ctx)
	defer close(d.playerChannel)
	var wg sync.WaitGroup
	go func() {
		for i := 0; i < d.config.MaxWorkers; i++ {
			wg.Add(1)
			go d.fetchData(&wg, done)

		}
		wg.Wait()
	}()
	for {
		select {
		case <-ctx.Done():
			if d.wantedTeams > 0 {
				return fmt.Errorf("there is %d items left", d.wantedTeams)
			}
			return nil
		}

	}

}

func (d *dataGetter) fetchData(wg *sync.WaitGroup, done context.CancelFunc) {
	defer wg.Done()
	for id := range d.idChannel {
		if d.wantedTeams <= 0 {
			done()
			return
		}
		var team model.Team
		url := fmt.Sprintf(d.config.Api, id)
		response, _ := http.Get(url)
		tools.Unmarshal(response, &team)
		if tools.Contains(d.config.WantedTeams, team.Data.Team.Name) {
			d.mu.Lock()
			d.wantedTeams--
			d.mu.Unlock()
			for _, player := range team.Data.Team.Players {
				player.Team = team.Data.Team.Name
				d.playerChannel <- player
			}
		}

	}
}
