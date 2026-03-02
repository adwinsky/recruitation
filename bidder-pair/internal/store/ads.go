package store

type Ad struct {
	ID string
}

func GetAds() []*Ad {
	// SQL Query
	return make([]*Ad, 0)
}
