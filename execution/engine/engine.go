package engine

import (
	"fmt"
	"github.com/stitchfix/flotilla-os/config"
	"github.com/stitchfix/flotilla-os/state"
)

//
// Engine defines the execution engine interface.
// An engine's sole purpose is to submit a run to
// backend resources
//
type Engine interface {
	Initialize(conf config.Config) error
	// v0
	Execute(definition state.Definition, run state.Run) (state.Run, error)

	// v1 - once runs contain a copy of relevant definition info
	// Execute(run state.Run) error

	Terminate(run state.Run) error
}

//
// NewExecutionEngine initializes and returns a new Engine
//
func NewExecutionEngine(conf config.Config) (Engine, error) {
	name := conf.GetString("execution_engine")
	if len(name) == 0 {
		name = "ecs"
	}

	switch name {
	case "ecs":
		eng := &ECSExecutionEngine{}
		return eng, eng.Initialize(conf)
	default:
		return nil, fmt.Errorf("No Engine named [%s] was found", name)
	}
}
