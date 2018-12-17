FROM golang AS builder
WORKDIR /go/src/portproxy
COPY . .
RUN go get -d .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

FROM scratch
WORKDIR /app
COPY --from=builder /go/src/portproxy/portproxy .
CMD [ "/app/portproxy" ]