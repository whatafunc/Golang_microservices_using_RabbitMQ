# Golang microservices using RabbitMQ

feat: producer & consumer microservices using separate config.yaml: take 3

Credits to github.com/OtusGolang/webinars_practical_part/29-queues


# install
```cd deployments && sudo docker compose up -d```

```cd ../```
```export DATABASE_URL="postgres://otus_user1:otus_password1@192.168.1.71:5433/events?sslmode=disable"```
```goose -dir migrations postgres $DATABASE_URL up```
```goose -dir migrations postgres $DATABASE_URL up```


# Usage
```go run cmd/calendar/*.go --config=configs/config.yaml```

```go run cmd/producer/*.go --config=configs/producer_config.yaml```

```go run cmd/consumer/*.go --config=configs/consumer_config.yaml```

# Github

```git add .
git commit -m "feat: xxx: take xx"
git push -u origin main```

# Run Producer Once a Day (via Cron / systemd timer)
```
0 7 * * * /usr/local/bin/producer --config=/etc/calendar/producer.yaml
``` 

### The Makefile provides several commands for building, running, and testing the Go application. Here's a breakdown:
   * `BIN`: Defines the output path for the compiled binary as ./bin calendar.                            
   * `DOCKER_IMG`: Sets the name for the Docker image to calendar:develop.                                
   * `GIT_HASH`: Retrieves the short Git commit hash of the current HEAD.                                 
   * `LDFLAGS`: Sets linker flags to embed build information (release, build date, Git hash) into the compiled binary.
   * `build`: Compiles the cmd/calendar Go application, including the LDFLAGS.
   * `run`: First builds the application, then runs it with the configuration from ./configs/config.yaml.
   * `build-img`: Builds a Docker image using build/Dockerfile, passing the LDFLAGS as a build argument and tagging the image as calendar:develop.
   * `run-img`: Builds the Docker image (if not already built) and then runs it.
   * `version`: Builds the application and then runs it with the version command-line argument to display its version information.
   * `test`: Runs Go tests for packages within ./internal/.... The commented-out line suggests it could also test ./pkg/.... 
   * `install-lint-deps`: Installs golangci-lint if it's not already present.
   * `lint`: Runs golangci-lint across all Go files in the project.
   * `.PHONY`: Declares all the above commands as phony targets, meaning they are not actual files.

  In essence, this Makefile streamlines common development tasks for this Go project, from building and running the application (both directly and in Docker) to testing and linting. 