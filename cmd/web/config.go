package web

type PackWebConfig struct {
	Packs []*PackEntry `json:"packs"`
}

type PackEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
