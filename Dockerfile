FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/chillaxio/chillax

ENV USER didip
ENV HTTP_ADDR 8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET sM-0YK3nK7m9XC2h

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://didip@localhost:5432/chillax?sslmode=disable

WORKDIR /go/src/github.com/chillaxio/chillax

RUN godep go build
RUN ./chillax

EXPOSE 8888