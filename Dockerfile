FROM golang:1.17.3-alpine as builder

WORKDIR /go/src/routinie-backend/
COPY go.mod go.sum ./
RUN apk add build-base
RUN go mod download
COPY . .
RUN GOOS=linux go build -a -o ./routinie-backend ./cmd/backend/main.go
RUN chmod +x ./routinie-backend

FROM alpine:latest
WORKDIR /
COPY --from=builder /go/src/routinie-backend/ .
CMD [ "./routinie-backend" ]