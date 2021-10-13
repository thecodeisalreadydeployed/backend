package preset

import (
	"bytes"
	"text/template"
)

type BuildOptions struct {
	BuildImage       string
	RunImage         string
	Executable       string
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

	text := presetText(framework)

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
FROM {{.BuildImage}} as build-env
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.InstallCommand}}
RUN {{.BuildCommand}}

FROM {{.RunImage}}
WORKDIR /app
COPY --from=build-env /app/{{.WorkingDirectory}}/package.json /app/{{.WorkingDirectory}}/yarn.lock ./
COPY --from=build-env /app/{{.WorkingDirectory}}/node_modules ./node_modules
COPY --from=build-env /app/{{.WorkingDirectory}}/{{.OutputDirectory}} ./{{.OutputDirectory}}
CMD {{.StartCommand}}	
`

	case FrameworkSpring:
		return `
FROM {{.BuildImage}} as build-env
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.BuildCommand}}

FROM {{.RunImage}}
WORKDIR /app
COPY --from=build-env /app/{{.WorkingDirectory}}/{{.OutputDirectory}}/{{.Executable}} .
CMD java -jar {{.Executable}}

`
	case FrameworkFlask:
		return `
FROM {{.RunImage}}
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.InstallCommand}}
CMD {{.StartCommand}}
`
	default:
		return `
FROM {{.RunImage}}
ADD . /app
WORKDIR /app/{{.WorkingDirectory}}
RUN {{.InstallCommand}}
RUN {{.BuildCommand}}
CMD {{.StartCommand}}	
`
	}
}
