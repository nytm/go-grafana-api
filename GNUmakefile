TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=github.com/nytm/go-grafana-api
GRAFANA_VERSION ?= latest
GO_VERSION ?= 1.10
GRAFANA_CONTAINER_NAME ?= grafana-dev
GRAFANA_ADMIN_PWD ?= pwd4test

default: test

test: start-grafana test-in-local

start-grafana:
	CONTAINER_COUNT=$$(docker ps -a -f "Name=$(GRAFANA_CONTAINER_NAME)" | wc -l) ;\
	if [[ $$CONTAINER_COUNT -lt 1 ]]; then \
		docker run --name $(GRAFANA_CONTAINER_NAME) -d -p localhost:3000:3000 -e "GF_SECURITY_ADMIN_PASSWORD=$(GRAFANA_ADMIN_PWD)" "grafana/grafana:$(GRAFANA_VERSION)"; \
	else \
		docker start $(GRAFANA_CONTAINER_NAME) ;\
	fi \

test-in-local:
	GRAFANA_AUTH=admin:$(GRAFANA_ADMIN_PWD) \
	go test \

test-in-docker:
	docker run -it -v `pwd`:/go/src/$(PKG_NAME) \
		--link $$GRAFANA_CONTAINER_ID:grafana \
		--workdir /go/src/$(PKG_NAME) \
		-e "GRAFANA_AUTH=admin:$(GRAFANA_ADMIN_PWD)" \
		-e "GRAFANA_URL=http://grafana:3000" "golang:$(GO_VERSION)" go test

fmt:
	gofmt -w $(GOFMT_FILES)
