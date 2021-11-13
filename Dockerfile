FROM golang:1.17.3-alpine As builder

WORKDIR /xyr/

RUN apk update && apk add musl-dev git gcc upx

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build --tags="json1 fts5 vacuum_incr preupdate_hook secure_delete_fast secure_delete" -ldflags "-s -w" -o /usr/bin/xyr

RUN upx -9 /usr/bin/xyr

FROM alpine

WORKDIR /xyr/

COPY --from=builder /usr/bin/xyr /usr/bin/xyr

CMD ["/usr/bin/xyr"]
