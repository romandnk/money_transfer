FROM golang:1.20 as build

WORKDIR /app

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o "./bin/app" ./cmd/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/

COPY ./config/config.yaml ./config/

CMD ["./bin/app"]