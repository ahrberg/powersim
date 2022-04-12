package powersim

import (
	"time"

	"github.com/robfig/cron/v3"
)

type CronConsumer struct {
	Power       int
	Sched       string
	Duration    time.Duration
	Description string
	CronShed    cron.Schedule
	on          bool
	onTime      time.Time
	offTime     time.Time
}

// NewCronConsumer return's a new cron consumer.
// The cron schedule specifies when the consumer is turned on.
// Cron schedule syntax "minute hour dayOfMonth month dayOfWeek"
// The duration determines how long the consumer in on and the power
// deterines the consumers constant power.
func NewCronConsumer(c CronConsumer) (*CronConsumer, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	s, err := parser.Parse(c.Sched)

	if err != nil {
		return nil, err
	}

	return &CronConsumer{
		CronShed:    s,
		Power:       c.Power,
		Sched:       c.Sched,
		Duration:    c.Duration,
		Description: c.Description,
	}, nil

}

func (s *CronConsumer) GetDescription() string {
	return s.Description
}

func (s *CronConsumer) GetPower(e Environment) int {

	if s.onTime.IsZero() {
		setTimes(s, e.time)
	}

	if s.on && e.time.Before(s.offTime) {
		return s.Power
	}

	if !s.on && (e.time.After(s.onTime) || e.time == s.onTime) {
		s.on = true
		return s.Power
	}

	if s.on && (e.time == s.offTime || e.time.After(s.offTime)) {
		s.on = false

		setTimes(s, e.time)
	}

	return 0
}

func setTimes(s *CronConsumer, t time.Time) error {
	next := s.CronShed.Next(t)
	s.onTime = next
	s.offTime = next.Add(s.Duration)

	return nil
}
