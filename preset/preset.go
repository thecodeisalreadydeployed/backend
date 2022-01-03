package preset

import (
	"bytes"
	"text/template"

	"github.com/thecodeisalreadydeployed/model"
)

type Framework string

const (
	FrameworkNestJS Framework = "FrameworkNestJS"
	FrameworkSpring Framework = "FrameworkSpring"
	FrameworkFlask  Framework = "FrameworkFlask"
	NoFramework     Framework = "NoFramework"
)

func Preset(opts model.BuildConfiguration, framework Framework) (string, error) {

	text := Text(framework)

	var buffer bytes.Buffer
	t := template.Must(template.New("Dockerfile").Parse(text))
	err := t.Execute(&buffer, opts)

	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func Text(framework Framework) string {
	switch framework {
	case FrameworkNestJS:
		return `
FROM node:14-alpine as build-env
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

	case FrameworkSpring:
		return `
FROM maven:3.8.3-ibmjava-8-alpine as build-env
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.BuildCommand}}

FROM openjdk:8-jdk-alpine
WORKDIR /app
COPY --from=build-env /app/{{.WorkingDirectory}}{{.OutputDirectory}}/*.jar .
CMD {{.StartCommand}}
`
	case FrameworkFlask:
		return `
FROM python:3.11.0a1-alpine3.14
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN pip install flask
RUN {{.InstallCommand}}
RUN {{.BuildCommand}}
CMD {{.StartCommand}}
`
	case NoFramework:
		return `
FROM alpine:latest
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.InstallCommand}}
RUN {{.BuildCommand}}
CMD {{.StartCommand}}	
`
	default:
		return ""
	}
}
