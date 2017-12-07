package worker

import (
	"fmt"
	"github.com/stitchfix/flotilla-os/config"
	"github.com/stitchfix/flotilla-os/execution/engine"
	flotillaLog "github.com/stitchfix/flotilla-os/log"
	"github.com/stitchfix/flotilla-os/state"
	"time"
)

//
// Worker defines a background worker process
//
type Worker interface {
	Initialize(
		conf config.Config, sm state.Manager, ee engine.Engine, log flotillaLog.Logger, pollInterval time.Duration) error
	Run()
}

func NewWorker(
	workerType string,
	log flotillaLog.Logger,
	conf config.Config,
	ee engine.Engine,
	sm state.Manager) (Worker, error) {

	var worker Worker

	switch workerType {
	case "submit":
		worker = &submitWorker{}
	case "retry":
		worker = &retryWorker{}
	case "status":
		worker = &statusWorker{}
	default:
		return nil, fmt.Errorf("No workerType %s exists", workerType)
	}

	pollIntervalString := conf.GetString(fmt.Sprintf("worker.%s_interval", workerType))
	if len(pollIntervalString) == 0 {
		return worker, fmt.Errorf("Worker type: [%s] needs worker.%s_interval set", workerType, workerType)
	}

	pollInterval, err := time.ParseDuration(pollIntervalString)
	if err != nil {
		return worker, err
	}
	err = worker.Initialize(conf, sm, ee, log, pollInterval)
	return worker, err
}
