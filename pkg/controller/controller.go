package controller

import (
	"fmt"
	"go.uber.org/zap"
	"locker/internal/config"
	"time"
)

func NewLockController(configuration *config.Configuration, log *zap.Logger) LockController {
	return &StandardLockController{
		configuration: configuration,
		log:           log,
	}
}

type LockController interface {
	Start() error
	Stop() error
}

type StandardLockController struct {
	running       bool
	log           *zap.Logger
	configuration *config.Configuration
}

func (lc *StandardLockController) Start() error {

	lc.running = true
	go func() {
		for lc.running {
			fmt.Println("ping...")
			time.Sleep(10 * time.Second)
		}
	}()

	return nil
}

func (lc *StandardLockController) Stop() error {

	lc.running = false
	return nil

}
