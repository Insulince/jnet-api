package trainingexecutor

import (
	"context"
	"fmt"
	"github.com/Insulince/jnet-api/untitled/cmd/api/networks"
	"github.com/Insulince/jnet-api/untitled/cmd/api/service"
	"github.com/Insulince/jnet-api/untitled/cmd/api/ws"
	"github.com/Insulince/jnet/pkg/trainer"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"io"
)

type (
	TrainingExecutor struct {
		service     service.Service
		workerCount int
	}

	TrainingWriter struct {
		NetworkId string
	}
)

var (
	_ io.Writer = TrainingWriter{}
)

func New(s service.Service, workerCount int) TrainingExecutor {
	var te TrainingExecutor

	te.service = s
	te.workerCount = workerCount

	return te
}

func (tw TrainingWriter) Write(data []byte) (int, error) {
	ws.NotifyActiveWorkers(tw.NetworkId, string(data))
	return len(data), nil
}

// TODO(justin): ERRORS are BROKEN in the worker pattern below.
func (te TrainingExecutor) Execute(ctx context.Context, tq networks.TrainingQueue, errs chan<- error) {
	defer close(errs)

	g, gctx := errgroup.WithContext(ctx)
	for i := 0; i < te.workerCount; i++ {
		func(id int) {
			g.Go(func() error {
				fmt.Printf("worker %v: ready for work\n", id)
				for ts := range tq {
					fmt.Printf("worker %v: starting training\n", id)

					var tw TrainingWriter
					tw.NetworkId = ts.NetworkId

					t := trainer.New(ts.TrainingConfiguration, ts.Data, tw)

					err := t.Train(ts.Network)
					if err != nil {
						return errors.Wrapf(err, "worker %v: training", id)
					}

					fmt.Printf("worker %v: done training, persisting changes\n", id)
					err = te.service.UpdateNetwork(gctx, ts.NetworkId, ts.Network)
					if err != nil {
						return errors.Wrapf(err, "worker %v: updating network", id)
					}

					fmt.Printf("worker %v: changes persisted\n", id)
				}

				return nil
			})
		}(i + 1)
	}

	err := g.Wait()
	if err != nil {
		errs <- errors.Wrap(err, "concurrently executing")
	}
}
