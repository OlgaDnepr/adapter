.PHONY: code-quality
code-quality:
	gometalinter --vendor --tests --exclude pb --exclude .*_gen\.go \
		--disable=gotype --disable=errcheck --disable=gas --disable=dupl \
		--deadline=1500s --checkstyle --sort=linter ./... > static-analysis.xml

.PHONY: dependencies
dependencies:
	glide install

.PHONY: mock
mock:
	mockgen -destination=mocks/mock_services_gen.go -package=mocks --source=pb/services.pb.go

.PHONY: proto
proto:
	protoc -I pb pb/services.proto --go_out=plugins=grpc:pb

.PHONY: generate-all
generate-all: proto mock

.PHONY: docker-build
docker-build: dependencies generate-all
	docker-compose up -d
	sleep 10
	docker start -i client

.PHONY: docker-clear
docker-clear:
	docker-compose stop
	docker-compose rm

.PHONY: test
test:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out