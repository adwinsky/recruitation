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

func (f *Filter) SelectAds(br *openrtb.BidRequest) ([]*store.Ad, bool, string) {
	now := time.Now()

	f.mu.Lock()
	defer f.mu.Unlock()

	// reset 1-minute window
	if now.Sub(f.window) >= time.Minute {
		f.window = now
		for k := range f.perUser {
			delete(f.perUser, k)
		}
	}

	if len(br.Imp) == 0 {
		return nil, false, "no_imp"
	}

	ads := store.GetAds()
	selected := make([]*store.Ad, 0, len(ads))

	// frequency capping per ad
	for _, ad := range ads {
		f.filterCounter++

		key := capKey{
			user: br.User.ID,
			ad:   ad.ID,
		}

		f.perUser[key]++
		if f.perUser[key] > f.cfg.MaxPerMinute {
			continue // skip only this ad
		}

		// expensive scoring simulation
		score := expensiveScore([]byte(br.User.ID), ad)
		if score%10 == 0 {
			continue
		}

		selected = append(selected, ad)
	}

	// random drop simulation (whole request)
	if f.rng.Intn(100) < 10 {
		return nil, false, "random_drop"
	}

	if len(selected) == 0 {
		return nil, false, "no_ads_after_filter"
	}

	return selected, true, "ok"
}
