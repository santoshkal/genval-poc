  FROM cgr.dev/chainguard/clang:latest as builder
  ENV APP_HOME=/app
  RUN apt-get update \
      && apt-get install -y gcc \
      && apt-get clean \
      && useradd -m -s /bin/bash -d $APP_HOME myappuser
  WORKDIR $APP_HOME
  COPY src/ $APP_HOME/src/ \
      && Makefile $APP_HOME/
  RUN make -C $APP_HOME

# STAGE 1
  FROM cgr.dev/chainguard/static:latest
  ENV APP_USER=myappuser
  ENV APP_HOME=/app
  WORKDIR $APP_HOME
  RUN useradd -m -s /bin/bash -d $APP_HOME $APP_USER
  COPY --from=builder $APP_HOME/src/myapp $APP_HOME/src/myapp
  USER $APP_USER
  CMD ["./src/myapp"]
