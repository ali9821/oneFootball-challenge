package main

import (
	"context"
	"fmt"
	"log"
	"project/cmd"
	"project/pkg/model"
	"time"
)

func main() {
	fmt.Println("task started")
	ctx := context.Background()
	startTime := time.Now()
	factory, err := cmd.NewFactory()
	if err != nil {
		log.Fatal(fmt.Sprintf("error in establishing factory : %d", err))
	}

	for _, runner := range factory.PipelineStages {
		go func(runner model.Runner) {
			err := runner.Run(ctx)
			if err != nil {
				log.Fatal(fmt.Sprintf("error in running : %d", err))
			}
		}(runner)
	}
	<-factory.Done
	fmt.Println(time.Since(startTime))
}
