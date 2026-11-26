package renderables

import (
	"encoding/json"

	"joshuanario.com/hyperfacts/pkg/hydratables"
)

type SuperRenderable struct {
	Catalog     `json:"Catalog,omitempty"`
	Renderables map[string]*Renderable
	Documents   map[string][]byte
}

func MarshalSuperRenderable(h *hydratables.Hydratable) ([]byte, error) {
	catalog, err := GetCatalog(h)
	if err != nil {
		return nil, err
	}
	renderableMap := make(map[string]*Renderable)
	for _, entityMap := range catalog.Networks {
		for _, slug := range entityMap {
			r, err := GetRenderable(slug, h)
			if err != nil {
				return nil, err
			}
			renderableMap[slug] = r
		}
	}
	documentMap := make(map[string][]byte)
	for _, document := range h.Folder.Documents {
		documentMap[document.Filename] = document.Bytes
	}

	return json.Marshal(SuperRenderable{
		Catalog:     *catalog,
		Renderables: renderableMap,
		Documents:   documentMap,
	})
}
