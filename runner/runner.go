package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/luizalabs/rey/checker"
	"github.com/luizalabs/rey/component"
	"github.com/luizalabs/rey/metric"
	"github.com/luizalabs/rey/notifier"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	cc     *checker.Checker
	ticker *time.Ticker
	nt     notifier.Notifier
}

func (r *Runner) Run(ctx context.Context, compList []*component.Component) error {
	if len(compList) == 0 {
		return fmt.Errorf("empty component list\n")
	}

	for range r.ticker.C {
		g, _ := errgroup.WithContext(ctx)
		for _, c := range compList {
			comp := c
			g.Go(func() error {
				st, err := r.cc.Check(comp)
				if err != nil {
					return err
				}
				if st.Status == comp.LastStatus {
					return nil
				}

				metric.NewGauge(st.Component).Set(float64(st.Status))
				if st.Status == checker.StatusDisruption {
					r.nt.Notify(
						fmt.Sprintf(
							"Component %s entering in Disruption status",
							st.Component,
						),
					)
				}

				comp.LastStatus = st.Status
				comp.LastDetail = st.Details
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) Stop() {
	r.ticker.Stop()
}

func New(interval int, cc *checker.Checker, nt notifier.Notifier) *Runner {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	return &Runner{cc: cc, ticker: ticker, nt: nt}
}
