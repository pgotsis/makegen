package main

var (
	buildFileTemplate = `#!/bin/bash
go build -o build/{{.Name}}`

	circleFileTemplate = `# Golang CircleCI 2.0 configuration file
#
# Check https://circleci.com/docs/2.0/language-go/ for more details
version: 2
jobs:
	build:
		docker:
			- image: circleci/golang:1.11
		working_directory: /go/src/github.com/ironarachne/{{.Name}}
		steps:
			- checkout
			- run: go get -v -t -d ./...
			- run: go test -v ./...
			- setup_remote_docker:
					docker_layer_caching: true
			- run: |
					TAG=0.1.$CIRCLE_BUILD_NUM
					docker build -t ironarachne/{{.Name}}:$TAG -t ironarachne/{{.Name}}:latest .
					docker login -u $DOCKER_USER -p $DOCKER_PASS
					docker push ironarachne/{{.Name}}:$TAG
					docker push ironarachne/{{.Name}}:latest
	deploy:
		machine:
				enabled: true
		steps:
			- run: curl -X POST 'https://portainer.ironarachne.com/api/webhooks/'

workflows:
	version: 2
	build-and-deploy:
		jobs:
			- build
			- deploy:
					requires:
						- build
					filters:
						branches:
							only: master`

	dockerFileTemplate = `# build stage
FROM golang:1.11 AS build-env
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/{{.Name}}

# final stage
FROM scratch
COPY --from=build-env /go/bin/{{.Name}} /go/bin/{{.Name}}
EXPOSE {{.Port}}
CMD ["/go/bin/{{.Name}}"]`

	gitignoreFileTemplate = `build/`

	mainFileTemplate = `package main

import (
	"github.com/ironarachne/utility"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		rand.Seed(time.Now().UnixNano())
		whatever := whatever()
		ctx.JSON(whatever)
	})

	app.Run(iris.Addr(":{{.Port}}"))
}`

	programFileTemplate = `package {{.Name}}

`

	readmeFileTemplate = `# {{.Name}}

Just another generator.`
)