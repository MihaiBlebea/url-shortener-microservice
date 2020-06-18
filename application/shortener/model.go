package shortener

import "time"

// Redirect is a model to represent an entity for shortener
type Redirect struct {
	ID      int       `json:"id"`
	URL     string    `json:"url"`
	Code    string    `json:"code"`
	Hits    int       `json:"hits"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func (r *Redirect) incrementHits() {
	r.Hits++
}
