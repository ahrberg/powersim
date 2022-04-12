# Power Simulation

Power simulation can simulate power consumers and calculate total energy consumption. For example model the power consumers in a house hold and calculate the consumed energy.

The simulation is written in Go and the simulation function `RunSim` return's a Go Channel that the caller can subscribe on. For every time step in the simulation the simulation state is reported on the Channel.

## Example

```go
// create one consumer running for 10 minutes
// every morning at 07.00 with power of 1000 W
c, err := NewCronConsumer(CronConsumer{
    Power:       1000,
    Duration:    time.Minute * 10,
    Sched:       "30 7 * * *",
    Description: "Hair dryer",
})

if err != nil {
    panic(err)
}

ctx := context.TODO()

// run the simulation from 2022-01-01 to 2022-01-02
ch := RunSim(ctx, Options{
    Consumers: []Consumer{c},
    StartTime: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
    EndTime:   time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC),
    Dt:        time.Second,
})

e := 0

// read the output from the simulation as it progresses
for r := range ch {
    // fmt.Println(r.T) // this will be the current time in the simulation
    // fmt.Println(r.P) // this will be the current power at time T in the simulation
    // fmt.Println(r.E) // this will be the accumulated energy consumption over the simulation
    e = r.E
}

fmt.Println(e)
// Output: 600000
```
