FROM golang:1.21 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .

# Ensure static linking and target platform compatibility
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go get -d -v ./...

# Build with additional flags for better compatibility
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o openapi .

FROM scratch AS runtime
COPY --from=build /go/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]


