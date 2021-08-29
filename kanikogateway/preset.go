package kanikogateway

import (
	"bytes"
	"text/template"
)

type BuildOptions struct {
	InstallCommand   string
	BuildCommand     string
	WorkingDirectory string
	OutputDirectory  string
	StartCommand     string
}

func PresetNestJS(opts BuildOptions) (string, error) {
	text := `
FROM node:14-alpine as build-env
WORKDIR /app
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.InstallCommand}}
RUN {{.BuildCommand}}

FROM node:14-alpine
WORKDIR /app
COPY --from=build-env /app/{{.WorkingDirectory}}package.json /app/{{.WorkingDirectory}}yarn.lock ./
COPY --from=build-env /app/{{.WorkingDirectory}}node_modules ./node_modules
COPY --from=build-env /app/{{.WorkingDirectory}}{{.OutputDirectory}} ./{{.OutputDirectory}}
CMD {{.StartCommand}}
	`

	var buffer bytes.Buffer
	t := template.Must(template.New("Dockerfile").Parse(text))
	err := t.Execute(&buffer, opts)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
