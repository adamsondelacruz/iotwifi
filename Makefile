IMAGE    ?= rpi-wifi
NAME     ?= rpiwifi
VERSION  ?= 1.0.0

all: build push

dev: dev_build dev_run

build:
	docker build -t $(IMAGE):latest .
	docker build -t $(IMAGE):arm32v6-$(VERSION) .

push:
	docker build -t $(IMAGE):latest .
	docker build -t $(IMAGE):arm32v6-$(VERSION) .

dev_build:
	docker build -t $(IMAGE) ./dev/

dev_run:
	sudo docker run --rm -it --privileged --network=host \
                   -v $(CURDIR):/go/src/github.com/cjimti/iotwifi \
                   -w /go/src/github.com/cjimti/iotwifi \
                   --name=$(NAME) $(IMAGE):latest


