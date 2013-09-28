export GOPATH=$(shell pwd)
install:
	@go install kindle-delivery
run:
	@pkill kindle-delivery || echo "no kindle-delivery process"
	@nohup ./bin/kindle-delivery 2>&1 >> ~/.logs/kindle-delivery/logs/kindle-delivery.log &
stop:
	@pkill kindle-delivery || echo "no kindle-delivery process"
get:
	@go get -u $(pkg)
test:
	@go install test
	@./bin/test
