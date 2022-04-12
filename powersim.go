package powersim

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type Environment struct {
	time time.Time
}

type Consumer interface {
	GetDescription() string
	GetPower(e Environment) int
}

// Options represents simulation options.
// Dt represents the time increment and resolution used in the simulation
// The simulation takes any Consumer that implements the Consumer interface
type Options struct {
	Consumers []Consumer
	StartTime time.Time
	EndTime   time.Time
	Dt        time.Duration
}

// Result contains P (power in Watt), E (energy in Joule) and T time.
type Result struct {
	P int
	E int
	T time.Time
}

// RunSim runs a simulation over the period of time specified in
// the options, ctx respects cancellation.
func RunSim(ctx context.Context, options Options) chan Result {
	t := options.StartTime
	eTot := 0 // Total energy over simulation
	resCh := make(chan Result, 100)

	go func() {
		defer close(resCh)
		for {
			if t == options.EndTime || t.After(options.EndTime) {
				break
			}

			pCh := runConsumers(options.Consumers, Environment{time: t})
			pTot := 0

			for pCh != nil {
				select {
				case p, ok := <-pCh:
					if !ok {
						pCh = nil
						continue
					}
					pTot += p

				case <-ctx.Done():
					return
				}
			}

			e := pTot * int(options.Dt.Seconds())
			eTot += e

			resCh <- Result{
				P: pTot,
				E: eTot,
				T: t,
			}
			t = t.Add(options.Dt)
		}
	}()

	return resCh
}

func runConsumers(consumers []Consumer, environment Environment) chan int {

	retCh := make(chan int)
	simCh := make(chan int, runtime.NumCPU()) // Simultaneous started power computations

	go func() {
		defer close(simCh)
		defer close(retCh)

		var wg sync.WaitGroup

		for _, c := range consumers {
			simCh <- 1
			wg.Add(1)
			go func(c Consumer) {
				p := c.GetPower(environment)
				retCh <- p
				<-simCh
				wg.Done()
			}(c)
		}

		wg.Wait()
	}()

	return retCh
}
