package renderables

import (
	"encoding/hex"
	"encoding/json"
	"hash/fnv"

	"joshuanario.com/hyperfacts/pkg/hydratables"
)

type Catalog struct {
	Subjects         []Subject
	RelationshipSets []RelationshipSet
	Networks         map[string]map[string]string
	Documents        []string
}

func MarshalCatalog(h *hydratables.Hydratable) ([]byte, error) {
	schemedEntities := sortedEntities(h)
	rsets := sortedRelationshipSets(h)
	subjects := make([]Subject, 0, len(schemedEntities))
	networks := map[string]map[string]string{}
	documents := []string{}
	if len(h.Documents) > 0 {
		for _, document := range h.Documents {
			documents = append(documents, document.Filename)
		}
	}
	for _, schemedEntity := range schemedEntities {
		entityStr := stringify(&schemedEntity)
		networks[entityStr] = make(map[string]string)
		for _, rset := range rsets {
			hash := hash(entityStr, rset.RoleURI, rset.Title)
			networks[entityStr][rset.RoleURI] = hash
		}
		name := schemedEntity.Scheme + "/" + schemedEntity.CharData
		hydratedName, err := hydratables.EntityQuery(schemedEntity.Scheme,
			schemedEntity.CharData)
		if err == nil {
			name = hydratedName
		}
		subjects = append(subjects, Subject{
			Name:   name,
			Entity: schemedEntity,
		})
	}
	return json.Marshal(Catalog{
		Subjects:         subjects,
		RelationshipSets: rsets,
		Networks:         networks,
		Documents:        documents,
	})
}

func hash(schemedEntity string, linkroleURI string, title string) string {
	h := fnv.New128a()
	h.Write([]byte(schemedEntity + linkroleURI + title))
	return hex.EncodeToString(h.Sum([]byte{}))
}
