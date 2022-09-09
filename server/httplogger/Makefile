build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o httplogger main.go

docker:
	docker build -t httplogger .
	docker save -o httplogger-docker.tar httplogger