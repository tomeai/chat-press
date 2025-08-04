IMAGE_VERSION := $(shell git rev-parse --short HEAD)

.PHONY: compile

clean:
	rm -rf ./releases

compile:clean
	@echo "Compile Project"
	go vet . && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.buildTime=`date +%Y%m%d.%H:%M:%S` -X main.buildCommit=`git rev-parse --short=12 HEAD` -X main.buildBranch=`git branch --show-current`" -o ./releases/seed-detect .

build:compile
	docker build --platform linux/amd64 . --file Dockerfile --tag registry.cn-beijing.aliyuncs.com/biyao/spider:$seed-detect-$(IMAGE_VERSION)

push:
	@echo "Pushing Docker image for $(MODULE) with version $(VERSION)"
	docker push registry.cn-beijing.aliyuncs.com/biyao/spider:$(MODULE)gateway-$(IMAGE_VERSION)