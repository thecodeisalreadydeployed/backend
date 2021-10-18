package preset

import (
	"github.com/thecodeisalreadydeployed/kanikogateway"
	"testing"
)

func TestNestJSPreset(t *testing.T) {
	actualText, err := Preset(kanikogateway.BuildOptions{
		InstallCommand:   "yarn install",
		BuildCommand:     "yarn build a",
		WorkingDirectory: "nx/",
		OutputDirectory:  "dist/apps/a",
		StartCommand:     "node dist/apps/a/main",
	}, FrameworkNestJS)

	if err != nil {
		t.Error(err)
	}

	expectedText := `
FROM node:14-alpine as build-env
ADD . /app
WORKDIR /app/nx/
RUN yarn install
RUN yarn build a

FROM node:14-alpine
WORKDIR /app
COPY --from=build-env /app/nx/package.json /app/nx/yarn.lock ./
COPY --from=build-env /app/nx/node_modules ./node_modules
COPY --from=build-env /app/nx/dist/apps/a ./dist/apps/a
CMD node dist/apps/a/main
`
	
	assert.Equal(t, actualText, expectedText)
}

func TestFlaskPreset(t *testing.T) {
	actualText, err := Preset(kanikogateway.BuildOptions{
		InstallCommand:   "echo",
		BuildCommand:     "echo",
		WorkingDirectory: "",
		OutputDirectory:  "",
		StartCommand:     "export FLASK_APP=app && flask run",
	}, FrameworkFlask)

	if err != nil {
		t.Error(err)
	}

	expectedText := `
FROM python:3.11.0a1-alpine3.14
ADD . /app
WORKDIR /app/
RUN pip install flask
RUN echo
RUN echo
CMD export FLASK_APP=app && flask run
`
	
	assert.Equal(t, actualText, expectedText)
}

func TestSpringPreset(t *testing.T) {
	actualText, err := Preset(kanikogateway.BuildOptions{
		InstallCommand:   "",
		BuildCommand:     "mvn clean install",
		WorkingDirectory: "src/main/java/com/example/demo/",
		OutputDirectory:  "target",
		StartCommand:     "java -jar *.jar",
	}, FrameworkSpring)

	if err != nil {
		t.Error(err)
	}

	expectedText := `
FROM maven:3.8.3-ibmjava-8-alpine as build-env
ADD . /app
WORKDIR /app/src/main/java/com/example/demo/
RUN mvn clean install

FROM openjdk:8-jdk-alpine
WORKDIR /app
COPY --from=build-env /app/src/main/java/com/example/demo/target/*.jar .
CMD java -jar *.jar
`
	
	assert.Equal(t, actualText, expectedText)
}
