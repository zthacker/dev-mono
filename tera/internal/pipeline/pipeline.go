package pipeline

import (
	"context"
	"fmt"
	"log"
	"tera/internal/protocols"
)

func RunPipeline(steps []protocols.PipelineStep, rawData []byte) {
	fmt.Printf("\n--- Starting Pipeline Input: %X ---\n", rawData)

	currentData := rawData
	var err error

	for _, step := range steps {
		fmt.Printf("[%s] Input: %X (%d bytes)\n", step.Name(), currentData, len(currentData))

		currentData, err = step.Process(context.Background(), currentData)
		if err != nil {
			log.Printf("[%s] FAILED: %v\n", step.Name(), err)
			return
		}
	}

	fmt.Printf("--- Pipeline Success! Final Output: %X ---\n", currentData)
}
