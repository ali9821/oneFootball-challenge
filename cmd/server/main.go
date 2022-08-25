package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"project/cmd"
	"project/pkg/model"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("task started")
	var mu sync.Mutex
	var wg sync.WaitGroup
	ctx := context.Background()

	// this is a comment for refactoring code below :
	// the issue is on every http call there are more open Goroutines that increase wait time to get the answer
	// and I should apply a fix that on evey response it kills them
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var players = make(map[string]model.PlayerData)
		if request.URL.Path != "/" {
			http.Error(writer, "%v", http.StatusNotFound)
		}
		factory, err := cmd.NewFactory()
		if err != nil {
			log.Fatal(fmt.Sprintf("error in establishing factory : %d", err))
		}
		start := time.Now()

		for _, runner := range factory.PipelineStages {
			go func(runner model.Runner) {
				err := runner.Run(ctx)
				if err != nil {
					log.Fatal(fmt.Sprintf("error in running : %d", err))
				}
			}(runner)
		}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for p := range factory.Response {
					mu.Lock()
					players[p.Id] = p
					mu.Unlock()
				}
			}()
		}
		wg.Wait()
		file, _ := json.MarshalIndent(players, "", " ")
		//_ = ioutil.WriteFile("test.json", file, 0644)
		fmt.Fprintf(writer, "%s", string(file))
		fmt.Println(runtime.NumGoroutine())
		fmt.Println(time.Since(start))
	})

	srv := &http.Server{
		Addr: ":8081",
	}
	func() {
		srv.ListenAndServe()
	}()

}
