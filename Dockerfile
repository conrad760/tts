FROM golang:alpine

RUN apk update \
  && apk add flite
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN	GOOS=linux go build -o main cmd/main.go
RUN adduser -S -D -H -h /app appuser
USER appuser
CMD ["./main"]
