package i18n

import (
	"github.com/tidwall/gjson"
	"io/fs"
)

type I18n struct {
	sources   map[string]interface{}
	sourceMap map[string]gjson.Result
	langName  string
}

func New(defaultLangName string) *I18n {
	i18n := &I18n{
		sources:   map[string]interface{}{},
		sourceMap: map[string]gjson.Result{},
		langName:  defaultLangName,
	}

	return i18n
}

func (a *I18n) SetDefaultLanguage(defaultLangName string) *I18n {
	a.langName = defaultLangName
	return a
}

func (a *I18n) GetDefaultLanguage() string {
	return a.langName
}

func (a *I18n) LoadFile(fsys fs.FS, path, langName string) (*I18n, error) {
	content, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	a.sources[langName] = string(content)
	a.sourceMap[langName] = gjson.ParseBytes(content)

	return a, nil
}

func (a *I18n) Parse(key string) string {
	value := a.sourceMap[a.langName].Get(key)
	if !value.Exists() {
		return ""
	}

	return value.String()
}

func (a *I18n) GetI18nSource() map[string]interface{} {
	return a.sources
}

func (a *I18n) GetLangName() string {
	return a.langName
}
