package asset_manifest

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"html/template"
)

const (
	script = `<script src="%s"></script>`
	style  = `<link href="%s" media="all" rel="stylesheet" />`
)

type AssetManifest struct {
	filename string
	assetRoot string
	manifest manifestStruct
}

type manifestStruct struct {
	Assets map[string]string
}

func NewAssetManifest(filename string, assetRoot string) (*AssetManifest, error) {
	manifest := AssetManifest{
		filename: filename,
		assetRoot: assetRoot,
	}

	return &manifest, manifest.Read()
}

func (m *AssetManifest) Read() error {
	file, err := ioutil.ReadFile(m.filename)

	if err != nil {
		return err
	}

	return json.Unmarshal(file, &m.manifest)
}

func (m *AssetManifest) AssetPath(asset string) string {
	return fmt.Sprintf("%s/%s", m.assetRoot, m.AssetName(asset))
}

func (m *AssetManifest) AssetName(asset string) string {
	hashed, ok := m.manifest.Assets[asset]

	if !ok {
		return asset
	} else {
		return hashed
	}
}

func (m *AssetManifest) JavascriptTag(path string) string {
	return fmt.Sprintf(script, m.AssetPath(path))
}

func (m *AssetManifest) StylesheetTag(path string) string {
	return fmt.Sprintf(style, m.AssetPath(path))
}

func (m *AssetManifest) GetHelpers() template.FuncMap {
	 return template.FuncMap {
		"javascript": helper(m.JavascriptTag),
		"stylesheet": helper(m.StylesheetTag),
		"asset_path": helper(m.AssetPath),
	}
}

func helper(fn func(string) string) (func(string) template.HTML) {
	return func(path string) template.HTML {
		return template.HTML(fn(path))
	}
}
