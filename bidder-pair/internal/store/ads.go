package store

type Ad struct {
	ID         string
	Categories []string
}

func GetAds() []*Ad {
	// SQL Query
	return make([]*Ad, 0)
}
