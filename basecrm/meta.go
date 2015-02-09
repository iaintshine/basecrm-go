package basecrm

type Meta struct {
	Type  string `json:"type"`
	Count int    `json:"count,omitempty"`
	Links *Links `json:"links,omitempty"`
}

type Links struct {
	First string `json:"first_page,omitempty"`
	Last  string `json:"last_page,omitempty"`
	Prev  string `json:"prev_page,omitempty"`
	Next  string `json:"next_page,omitempty"`
}
