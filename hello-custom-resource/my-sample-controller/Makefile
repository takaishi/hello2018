date := $(shell date +'%s')

default:
	env GO111MODULE=off GOOS=linux GOARCH=amd64 go build -o controller-main main.go
	docker build . -t rtakaishi/sample-controller-main
	kubectl patch deploy controller-main -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"$(date)\"}}}}}"