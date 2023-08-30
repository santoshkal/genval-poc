  FROM cgr.dev/chainguard/go:latest As builder
  ENV APP_HOME=/app
  RUN useradd -m -s /bin/bash -d $APP_HOME myappuser
  WORKDIR $APP_HOME
  COPY go.mod go.sum $APP_HOME/
  RUN apt-get update \
      && apt-get clean \
      && go mod download
  COPY src/ $APP_HOME/src/
  RUN CGO_ENABLED=0 go build -o myapp $APP_HOME/src/main.go

# STAGE 1
  FROM cgr.dev/chainguard/static:latest
  ENV APP_USER=myappuser
  ENV APP_HOME=/app
  RUN useradd -m -s /bin/bash -d $APP_HOME $APP_USER
  WORKDIR $APP_HOME
  COPY --from=builder $APP_HOME/myapp $APP_HOME/myapp
  EXPOSE 8080
  USER $APP_USER
  CMD ["./myapp"]
