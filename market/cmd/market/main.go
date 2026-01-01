package main

import (
	"log"
	"market/internal/services/engine"
)

func main() {
	mEng := engine.NewEngine(engine.GameEngineConfig{
		MapName: "Arena",
		Mode:    engine.ModeMultiPlayer,
	})

	if err := mEng.Initialize(); err != nil {
		panic(err)
	}

	if err := mEng.Start(); err != nil {
		log.Println("error running game:", err)
	}

}
