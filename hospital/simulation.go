package hospital

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// SimConfig configures a concurrent simulation run by Hospital.Simulate.
// It is a plain configuration value, not an entity: it has no invariants to
// protect between calls, so its fields are exported and Simulate validates
// them before using them.
type SimConfig struct {
	// Rounds is how many times every patient falls asleep and wakes up.
	// It must be positive.
	Rounds int
	// MaxDelay is the longest random pause before a sleep attack and
	// before waking up. It must not be negative.
	MaxDelay time.Duration
	// Locations are the places where a sleep attack can happen; each attack
	// picks one at random. There must be at least one, and none blank.
	Locations []string
}

// SimResult counts what happened during a simulation.
type SimResult struct {
	// Attacks is how many sleep attacks were registered in the history,
	// including the ones where the patient stayed in the hallway.
	Attacks int
	// RoomFailures is how many of those Attacks found no free room
	// (ErrNoRoomAvailable). They are not counted as Errors.
	RoomFailures int
	// Errors is how many calls to HandleSleepAttack or WakeUpPatient failed
	// for any other reason.
	Errors int
}

// Simulate runs one goroutine per admitted patient. Each goroutine repeats
// cfg.Rounds times: wait a random time up to cfg.MaxDelay, suffer a sleep
// attack at a random location of cfg.Locations (HandleSleepAttack), wait
// another random time and wake up (WakeUpPatient). Simulate waits for every
// goroutine to finish and returns the counters. It returns an error wrapping
// ErrInvalidData if cfg is invalid; failures inside the goroutines are
// counted in the result, not returned.
//
// The goroutines only touch the hospital through its exported methods,
// which take the mutex. Simulate itself never holds the mutex, so it must
// not be called from a ...Locked method.
func (h *Hospital) Simulate(cfg SimConfig) (SimResult, error) {
	if err := cfg.validate(); err != nil {
		return SimResult{}, fmt.Errorf("simulate: %w", err)
	}

	// Several goroutines increase these counters at the same time, so they
	// are atomic: a plain n++ is a read followed by a write, and two
	// goroutines could read the same value and lose one increment.
	var attacks, roomFailures, errCount atomic.Int64

	var wg sync.WaitGroup
	for _, p := range h.Patients() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range cfg.Rounds {
				time.Sleep(randomDelay(cfg.MaxDelay))
				location := cfg.Locations[rand.Intn(len(cfg.Locations))]
				_, err := h.HandleSleepAttack(p, location)
				switch {
				case err == nil:
					attacks.Add(1)
				case errors.Is(err, ErrNoRoomAvailable):
					// The episode was registered: the patient sleeps in the hallway.
					attacks.Add(1)
					roomFailures.Add(1)
				default:
					errCount.Add(1)
				}

				time.Sleep(randomDelay(cfg.MaxDelay))
				if err := h.WakeUpPatient(p); err != nil {
					errCount.Add(1)
				}
			}
		}()
	}
	wg.Wait()

	return SimResult{
		Attacks:      int(attacks.Load()),
		RoomFailures: int(roomFailures.Load()),
		Errors:       int(errCount.Load()),
	}, nil
}

// validate checks the rules documented in SimConfig. It returns an error
// wrapping ErrInvalidData for the first rule that is broken.
func (cfg SimConfig) validate() error {
	if cfg.Rounds <= 0 {
		return fmt.Errorf("rounds must be positive, got %d: %w", cfg.Rounds, ErrInvalidData)
	}
	if cfg.MaxDelay < 0 {
		return fmt.Errorf("max delay must not be negative, got %v: %w", cfg.MaxDelay, ErrInvalidData)
	}
	if len(cfg.Locations) == 0 {
		return fmt.Errorf("no locations: %w", ErrInvalidData)
	}
	for i, loc := range cfg.Locations {
		if strings.TrimSpace(loc) == "" {
			return fmt.Errorf("location %d is blank: %w", i, ErrInvalidData)
		}
	}
	return nil
}

// randomDelay returns a random duration between 0 and limit, both included.
// The top-level functions of math/rand are safe for concurrent use.
func randomDelay(limit time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(limit) + 1))
}
