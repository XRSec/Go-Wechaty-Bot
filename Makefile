.PHONY: all
all:
	@reset;bash Gateway/wechatyGateway.sh

.PHONY: server
server:
	@reset; cd Server; go run main.go

.PHONY: test
test:
	@reset;cd Server; air
