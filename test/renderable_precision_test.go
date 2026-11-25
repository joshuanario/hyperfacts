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

func TestMarshalRenderable_High_Precision(t *testing.T) {
	hcache := gocache.New(gocache.NoExpiration, gocache.NoExpiration)
	serializables.WorkingDirectoryPath = filepath.Join(".", "wd")
	serializables.GlobalTaxonomySetPath = filepath.Join(".", "gts")
	hydratables.InjectCache(hcache)
	zipFile := filepath.Join(".", "wd", "high_precision.zip")
	workingDir := filepath.Join(serializables.WorkingDirectoryPath, "folders", "high_precision")
	defer func() {
		os.RemoveAll(workingDir)
	}()
	err := unZipTestData(workingDir, zipFile)
	if err != nil {
		panic("Error: " + err.Error())
	}
	_, err = os.Stat(workingDir)
	if os.IsNotExist(err) {
		t.Fatalf("%s", "Error: "+err.Error())
		return
	}
	f, err := serializables.Discover("high_precision")
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	h, err := hydratables.Hydrate(f)
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	slug := "b8d6d55ca3d34f8be977a0cab0bae837"
	data, err := renderables.MarshalRenderable(slug, h)
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	r := renderables.Renderable{}
	err = json.Unmarshal(data, &r)
	if err != nil {
		t.Fatalf("%s", "Error: "+err.Error())
	}
	if r.RelationshipSet.RoleURI != "http://www.workiva.com/role/StockBasedCompensationEmployeeStockPurchasePlanDetails" {
		t.Fatalf("expected http://www.workiva.com/role/StockBasedCompensationEmployeeStockPurchasePlanDetails; outcome %s;\n", r.RelationshipSet.RoleURI)
	}
}
