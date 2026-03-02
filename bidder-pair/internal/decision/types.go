package decision

import "bidder-pair/internal/store"

type DecisionRequest struct {
	RequestID string      `json:"request_id"`
	UserID    string      `json:"user_id"`
	ImpIDs    []string    `json:"imp_ids"`
	Ads       []*store.Ad `json:"ads"`
}

type DecisionResponse struct {
	Allow bool    `json:"allow"`
	Price float64 `json:"price"`
	AdID  string  `json:"ad_id"`
}
