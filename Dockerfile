FROM golang:alpine3.11 AS build

RUN apk --no-cache add build-base git bzr mercurial gcc ca-certificates tzdata

WORKDIR /src

ADD src/go.mod /src
ADD src/go.sum /src

RUN go mod download

ADD src/*.go /src/
ADD src/bot/*.go /src/bot/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gobot .

FROM scratch

COPY --from=build /src/gobot /bot
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo

CMD ["/bot"]
