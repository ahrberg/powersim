package powersim

import (
	"context"
	"testing"
	"time"
)

func TestSim(t *testing.T) {

	cSpec := []CronConsumer{
		{
			Power:       1000,
			Duration:    time.Minute * 10,
			Sched:       "30 7 * * *",
			Description: "Hair dryer",
		},
	}

	consumers := []Consumer{}

	for _, cs := range cSpec {
		c, err := NewCronConsumer(cs)

		if err != nil {
			panic(err)
		}

		consumers = append(consumers, c)
	}

	ctx := context.TODO()

	ch := RunSim(ctx, Options{
		Consumers: consumers,
		StartTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
		Dt:        time.Second,
	})

	e := 0
	eExp := 600000
	calls := 0
	callsExp := 3600 * 24

	for r := range ch {
		calls += 1
		e = r.E
	}

	if e != eExp {
		t.Errorf("Expected total energy to be %d got %d", eExp, e)
	}

	if calls != callsExp {
		t.Errorf("Expected %d calls got %d calls", callsExp, calls)
	}
}
