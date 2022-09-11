build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o httplogger main.go

deploy:
	docker compose up -d
