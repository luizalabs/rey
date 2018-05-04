package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/luizalabs/rey/checker"
	"github.com/luizalabs/rey/component"
	"golang.org/x/sync/errgroup"
)

type Runner struct {
	cc     *checker.Checker
	ticker *time.Ticker
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
				if st.StatusID == comp.LastStatus {
					return nil
				}

				comp.LastStatus = st.StatusID
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

func New(interval int, cc *checker.Checker) *Runner {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	return &Runner{cc: cc, ticker: ticker}
}
