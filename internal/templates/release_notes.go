package templates

import (
	"strings"
	"text/template"
	"time"
)

func RenderHeader(start, end time.Time) (string, error) {

	// message for start of thread
	headerTpl, err := template.New("header").Parse(MustAssetString("templates/header.md.gotmpl"))
	if err != nil {
		return "", err
	}
	data := map[string]any{
		"start": start.Format(time.RFC822),
		"end":   end.Format(time.RFC822),
	}
	var headerBuilder strings.Builder
	if err = headerTpl.Execute(&headerBuilder, data); err != nil {
		return "", err

	}
	return headerBuilder.String(), nil
}
