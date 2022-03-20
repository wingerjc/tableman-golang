package web

type PackWebConfig struct {
	Packs []*PackEntry `json:"packs"`
}

type PackEntry struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Path  string `json:"path"`
}
