package hyperfacts_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	gocache "github.com/patrickmn/go-cache"
	"joshuanario.com/hyperfacts/pkg/hydratables"
	"joshuanario.com/hyperfacts/pkg/renderables"
	"joshuanario.com/hyperfacts/pkg/serializables"
)

func TestMarshalRenderable_Erroneous_Footnote(t *testing.T) {
	hcache := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	serializables.WorkingDirectoryPath = filepath.Join(".", "wd")
	serializables.GlobalTaxonomySetPath = filepath.Join(".", "gts")
	hydratables.InjectCache(hcache)
	workingDir := filepath.Join(serializables.WorkingDirectoryPath, "folders", "test_erroneous")
	_, err := os.Stat(workingDir)
	if os.IsNotExist(err) {
		t.Fatalf("%s", "Error: "+err.Error())
		return
	}
	f, err := serializables.Discover("test_erroneous")
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	h, err := hydratables.Hydrate(f)
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	slug := "d04179d4c5776ebc8842ee1fbb20f091"
	data, err := renderables.MarshalRenderable(slug, h)
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	r := renderables.Renderable{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	if r.RelationshipSet.RoleURI != "http://nationalbeverage.com/20200502/role/statement-note-10-stockbased-compensation-summary-of-stock-option-activity-details" {
		t.Fatalf("expected http://nationalbeverage.com/20200502/role/statement-note-10-stockbased-compensation-summary-of-stock-option-activity-details; outcome %s;\n", r.RelationshipSet.RoleURI)
	}
	if (*r.PGrid.FactualQuadrant[6][2])["Unlabelled"].Core != "11.14" {
		t.Fatalf("expected 11.14; outcome %s;\n", (*r.PGrid.FactualQuadrant[6][2])["Unlabelled"].Core)
	}
	if len(r.PGrid.FootnoteGrid[6][2]) != 1 {
		t.Fatalf("expected 1 footnote link; outcome %d;\n", len(r.PGrid.FootnoteGrid[6][2]))
	}
	if len(r.PGrid.Footnotes) != 1 {
		t.Fatalf("expected 1 footnote; outcome %d;\n", len(r.PGrid.Footnotes))
	}
	if r.PGrid.Footnotes[0] != "Weighted average exercise price." {
		t.Fatalf("expected Weighted average exercise price.; outcome %s;\n", r.PGrid.Footnotes[0])
	}
}
