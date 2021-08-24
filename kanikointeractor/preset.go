package kanikointeractor

import "fmt"

func PresetNestJS(installCommand string, buildCommand string, outputDirectory string, startCommand string) string {
	template := `
FROM node:14-alpine as build-env
WORKDIR /app
ADD . /app
RUN %s
RUN %s

FROM node:14-alpine
WORKDIR /app
COPY --from=build-env /app/package.json /app/yarn.lock ./
COPY --from=build-env /app/node_modules ./node_modules
COPY --from=build-env /app/%s ./%s
CMD %s
	`

	return fmt.Sprintf(template, installCommand, buildCommand, outputDirectory, outputDirectory, startCommand)
}
