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

	for _, ad := range ads {
		f.filterCounter++

		key := capKey{
			user: br.User.ID,
			ad:   ad.ID,
		}

		// aerospike request?
		f.perUser[key]++
		if f.perUser[key] > f.cfg.MaxPerMinute {
			continue
		}

		// repeatedSegmentCheck has high CPU cost and rejects ~30% of ads on average.
		if repeatedSegmentCheck(br.User.ID) {
			continue
		}

		// requestHashFilter has low CPU cost and rejects 70% of ads on average.
		if requestHashFilter(ad, br)%5 == 0 {
			continue
		}

		// categoryFilter has high CPU cost It rejects ~50% of ads on average.
		// br.Site.Category has low cardinality
		if !siteCategoryFilter(ad, br.Site.Category) {
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

// repeatedSegmentCheck has high CPU cost and rejects ~30% of ads on average.
func repeatedSegmentCheck(userID string) bool {
	time.Sleep(150 * time.Microsecond)
	return len(userID)%3 == 0
}

// requestHashFilter has low CPU cost and rejects 70% of ads on average.
func requestHashFilter(ad *store.Ad, br *openrtb.BidRequest) int {
	time.Sleep(100 * time.Microsecond)
	return len(br.User.ID) + len(ad.ID)
}

func siteCategoryFilter(ad *store.Ad, category string) bool {
	for _, c := range ad.Categories {
		if c == category {
			return true
		}
	}
	return false
}
