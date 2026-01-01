package engine

import (
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"time"
)

type Status int
type Mode string

const (
	ModeSinglePlayer Mode = "single_player"
	ModeMultiPlayer  Mode = "multi_player"
)

const (
	StatusRunning      Status = 0
	StatusStopped      Status = 1
	StatusPaused       Status = 2
	StatusError        Status = 3
	StatusUnknown      Status = 4
	StatusInitializing Status = 5
	StatusEnded        Status = 6
)

type Player struct {
	ID    int
	Name  string
	Score int
}

type GameEngineConfig struct {
	MaxPlayers    int
	CurrentPlayer int
	Players       []Player
	MapName       string
	Mode          Mode
	Difficulty    string
}

type GameEngine interface {
	// Initialize sets up the game engine with necessary configurations.
	Initialize() error
	// Start initiates the game engine.
	Start() error
	// Stop terminates the game engine.
	Stop() error
	// Status returns the current status of the game engine.
	Status() Status
}

type multiPlayer struct {
	config GameEngineConfig
	status Status
}

func NewEngine(config GameEngineConfig) GameEngine {
	return &multiPlayer{
		config: config,
		status: StatusInitializing,
	}
}

func (e *multiPlayer) Initialize() error {
	log.Println("initializing multiplayer game engine")
	e.config.MaxPlayers = rand.IntN(16)
	if e.config.MaxPlayers <= 2 {
		e.config.MaxPlayers = 4
	}
	
	log.Println("max players set to:", e.config.MaxPlayers)
	e.config.Players = make([]Player, e.config.MaxPlayers)
	// Simulate loading players
	for i := 0; i < e.config.MaxPlayers; i++ {
		player := Player{
			ID:    i + 1,
			Name:  fmt.Sprintf("Player%d", i+1),
			Score: 0,
		}
		e.config.Players = append(e.config.Players, player)
	}
	fmt.Printf("loaded players: %+v\n", e.config.Players)

	e.config.CurrentPlayer = 0

	e.status = StatusRunning
	return nil
}

func (e *multiPlayer) Start() error {

	for e.status != StatusEnded {
		switch e.status {
		case StatusPaused:
			log.Println("game is paused")
			time.Sleep(1 * time.Second)
			e.status = StatusRunning
		case StatusError:
			return errors.New("game encountered an error, attempting to recover")
		case StatusStopped:
			// Handle stopped state.
			log.Println("game was stopped, reinitializing")
			time.Sleep(2 * time.Second)
			e.Initialize()
		case StatusRunning:
			// Main game loop logic goes here.
			player := e.nextPlayer()
			log.Printf("player %d is taking action\n", player.ID)
			time.Sleep(500 * time.Millisecond)
			if player.ID%e.config.MaxPlayers == 0 {
				log.Println("player is pausing the game")
				e.status = StatusPaused
			}
			if player.ID%e.config.MaxPlayers == 1 {
				log.Println("player is throwing an error")
				e.status = StatusError
			}
			if player.ID%e.config.MaxPlayers == 2 {
				log.Println("player is stopping the game")
				e.status = StatusStopped
			}
		default:
			log.Println("unknown status, ending game")
			e.status = StatusEnded
		}
	}

	return nil
}

func (e *multiPlayer) Stop() error {
	e.status = StatusStopped
	return nil
}

func (e *multiPlayer) Status() Status {
	return e.status
}

func (e *multiPlayer) statustoString() string {
	switch e.status {
	case StatusRunning:
		return "running"
	case StatusStopped:
		return "stopped"
	case StatusPaused:
		return "paused"
	case StatusError:
		return "error"
	case StatusInitializing:
		return "initializing"
	default:
		return "unknown"
	}
}

func (e *multiPlayer) nextPlayer() Player {
	if len(e.config.Players) == 0 {
		return Player{}
	}
	index := rand.IntN(len(e.config.Players))
	return e.config.Players[index]
}
