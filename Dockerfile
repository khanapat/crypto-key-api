FROM golang:alpine as build

WORKDIR /app

ADD . /app

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go build -o goapp -ldflags '-w -s' main.go

# ---

FROM scratch as final

WORKDIR /app

COPY --from=build /app/goapp ./
COPY --from=build /app/config/config.yaml ./config/
ENV TZ=Asia/Bangkok

ENTRYPOINT ["/app/goapp"]