all: bin/xkcdbot
	docker build -t faryon93/xkcdbot:latest .

bin/xkcdbot: main.go
	CGO_ENABLED=0 bowler build