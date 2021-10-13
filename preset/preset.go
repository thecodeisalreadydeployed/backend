package preset

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

type Framework string

const (
	FrameworkNestJS Framework = "FrameworkNestJS"
	FrameworkSpring Framework = "FrameworkSpring"
	FrameworkFlask  Framework = "FrameworkFlask"
)

func Preset(opts BuildOptions, framework Framework) (string, error) {

	text := "FROM " + image(framework) + ` as build-env
WORKDIR /app
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.InstallCommand}}
RUN {{.BuildCommand}}

FROM ` + image(framework) + `
WORKDIR /app` + presetText(framework) + `COPY --from=build-env /app/{{.WorkingDirectory}}{{.OutputDirectory}} ./{{.OutputDirectory}}
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

func presetText(framework Framework) string {
	switch framework {
	case FrameworkNestJS:
		return `
COPY --from=build-env /app/{{.WorkingDirectory}}package.json /app/{{.WorkingDirectory}}yarn.lock ./
COPY --from=build-env /app/{{.WorkingDirectory}}node_modules ./node_modules
`
	case FrameworkSpring:
		return `
COPY --from=build-env /app/{{.WorkingDirectory}}pom.xml /app/{{.WorkingDirectory}}*gradle* /app/{{.WorkingDirectory}}*mvn* ./
COPY --from=build-env /app/{{.WorkingDirectory}}target ./target
`
	case FrameworkFlask:
		return "\n"
	default:
		return "\n"
	}
}

func image(framework Framework) string {
	switch framework {
	case FrameworkNestJS:
		return "node:14-alpine"
	case FrameworkSpring:
		return "openjdk:8-jdk-alpine"
	case FrameworkFlask:
		return "3.11.0a1-alpine3.14"
	default:
		return "alpine:latest"
	}
}
