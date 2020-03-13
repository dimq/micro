.PHONY: all
suffix=v
buildName=micro
all: go-build docker-build docker-push

 

versionTarget:
ifeq ($(strip $(version)),)
myversion := latest
else
myversion := $(suffix)$(version)
endif

go-build:
	go build -o $(buildName)

docker-build: versionTarget go-build
	docker build -t $(buildName):$(myversion) .    

 

docker-push:
	docker image push $(buildName):$(myversion)
