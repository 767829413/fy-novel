package chapter

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"

	"fy-novel/internal/definition"
)

var templates sync.Map

func templateConvert(title, content, extName string) (string, error) {
	tmpl, err := getTemplate(extName)
	if err != nil {
		return "", err
	}

	data := struct {
		Title   string
		Content string
	}{
		Title:   title,
		Content: content,
	}

	var renderedContent bytes.Buffer
	err = tmpl.Execute(&renderedContent, data)
	if err != nil {
		return "", err
	}

	return renderedContent.String(), nil
}

func getTemplate(extName string) (*template.Template, error) {
	if tmpl, ok := templates.Load(extName); ok {
		return tmpl.(*template.Template), nil
	}

	var templateContent string
	switch extName {
	case definition.NovelExtname_HTML: // Generic HTML Templates
		templateContent = definition.NovelTemp_HTML
	case definition.NovelExtname_EPUB:
		templateContent = definition.NovelTemp_EPUB
	default:
		return nil, fmt.Errorf("unsupported template type: %s", extName)
	}

	tmpl, err := template.New(extName).Parse(templateContent)
	if err != nil {
		return nil, err
	}

	actual, loaded := templates.LoadOrStore(extName, tmpl)
	if loaded {
		// Another goroutine has already stored a value, use that
		return actual.(*template.Template), nil
	}

	return tmpl, nil
}
