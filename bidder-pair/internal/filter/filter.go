package filter

import (
	"math/rand"
	"sync"
	"time"

	"bidder-pair/internal/openrtb"
	"bidder-pair/internal/store"
)

type Config struct {
	MaxPerMinute int
}

// frequency cap key: user + ad
type capKey struct {
	user string
	ad   string
}

type Filter struct {
	cfg Config

	mu            sync.Mutex
	perUser       map[capKey]int
	window        time.Time
	rng           *rand.Rand
	filterCounter int
}

func New(cfg Config) *Filter {
	return &Filter{
		cfg:     cfg,
		perUser: make(map[capKey]int),
		window:  time.Now(),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (f *Filter) Filter(br *openrtb.BidRequest) (bool, string) {
	now := time.Now()

	// reset 1-minute window
	if now.Sub(f.window) >= time.Minute {
		f.window = now
		for k := range f.perUser {
			delete(f.perUser, k)
		}
	}

	// frequency capping per ad
	for _, ad := range store.GetAds() {
		f.filterCounter++

		key := capKey{
			user: br.User.ID,
			ad:   ad.ID,
		}

		f.perUser[key]++

		if f.perUser[key] > f.cfg.MaxPerMinute {
			return false, "freqcap"
		}

		// expensive scoring simulation
		score := expensiveScore([]byte(br.User.ID), ad)
		if score%10 == 0 {
			return false, "low_score"
		}
	}

	// random drop simulation
	drop := f.rng.Intn(100) < 10
	if drop {
		return false, "random_drop"
	}

	if len(br.Imp) == 0 {
		return false, "no_imp"
	}

	return true, "ok"
}
