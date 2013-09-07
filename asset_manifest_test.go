package asset_manifest

import (
	"testing"
	"html/template"
	"bytes"
	"strings"
)

const (
	fixture    = "fixtures/manifest.json"
	assetRoot  = "/assets"
	appjs      = "app-a9b6dec34dcf660f443f820e37a77958.js"
	appjspath  = "/assets/app-a9b6dec34dcf660f443f820e37a77958.js"
    appcsspath = "/assets/app-ce73112ab9cf7e1091b2275e8062d37c.css"
)

func Test_NewAssetManifestWithBadManifestPathReturnError(t *testing.T) {
	_, err := NewAssetManifest("badpath", "")

	if err == nil {
		t.Error("Expected a bad manifest file to return an error")
	}
}

func Test_NewAssetManifestWithGoodManifestPathDoesNotReturnError(t *testing.T) {
	_, err := NewAssetManifest(fixture, assetRoot)

	if err != nil {
		t.Error(err)
	}
}

func manifest(t *testing.T) (*AssetManifest) {
	manifest, err := NewAssetManifest(fixture, assetRoot)
	if err != nil {
		t.Error(err)
	}

	return manifest
}

func Test_AssetPathWithNotFoundAsset(t *testing.T) {
	manifest := manifest(t)

	path := manifest.AssetPath("notfound.js")

	if path != "/assets/notfound.js" {
		t.Errorf("expected: /assets/notfound.js got %s", path)
	}
}

func Test_AssetPathReturnsFullAssetPath(t *testing.T) {
	manifest := manifest(t)

	path := manifest.AssetPath("app.js")

	if path != appjspath {
		t.Errorf("expected: %s got %s", appjspath, path)
	}
}

func Test_JavascriptTag(t *testing.T) {
	manifest := manifest(t)

	tag := manifest.JavascriptTag("app.js")

	if !strings.Contains(tag, appjspath) {
		t.Error("expected script tag with hashed asset path, got: %s", tag)
	}
}

func Test_StylesheetTag(t *testing.T) {
	manifest := manifest(t)

	tag := manifest.StylesheetTag("app.css")

	if !strings.Contains(tag, appcsspath) {
		t.Error("expected script tag with hashed asset path, got: %s", tag)
	}
}

func Test_HelperFunctions(t *testing.T) {
	manifest := manifest(t)

	helpers := manifest.GetHelpers()

	tpl := template.New("foo").Funcs(helpers)
	tpl, err := tpl.Parse(`{{asset_path "app.js"}}`)

	if err != nil {
		t.Error(err)
	}

	buffer := new(bytes.Buffer)
	tpl.Execute(buffer, nil)

	output := string(buffer.Bytes())

	if output != appjspath {
		t.Errorf("expected: %s got %s", appjspath, output)
	}
}

