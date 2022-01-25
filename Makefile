build:
	GOOS=linux go build cmd/main.go
	docker build --platform "linux/amd64" -t gcr.io/text-to-speech-338303/say .
	rm -f main

push:
	docker push gcr.io/text-to-speech-338303/say
