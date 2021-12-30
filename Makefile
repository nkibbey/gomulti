LEN := 0
IMG_REPO := gomulti
DOCKER_REPO := nkibbey/gomulti

ifeq ($(shell test -d .git; echo $$?), 0)
  ifeq ($(origin CI_COMMIT_SHA), undefined)
    CI_COMMIT_SHA := $(shell git describe --all --dirty --long)
  endif

  LVER := $(shell git describe --tags)
  SVER := $(shell git describe --tags | rev | cut -c 10- | rev)

  ifeq ($(shell test -z $(LVER); echo $$?), 1)
    LEN := $(shell expr length $(LVER))
  endif

  ifeq ($(shell test $(LEN) -gt 10; echo $$?), 0)
 	VER := $(SVER)
  else
	VER := $(LVER)
  endif
endif

define DOCKFILE
FROM golang:1.17-alpine
COPY bin/* /go/bin/
ENTRYPOINT ["/go/bin/goreel"]
endef
export DOCKFILE

build:
	@mkdir -p bin
	go build -ldflags="-X 'main.Version=$(VER)' -X 'main.GitCommit=$(CI_COMMIT_SHA)' -X 'main.BuildTime=$$(date)'" -tags netgo -o bin ./...

oci: build
	@echo -e "$$DOCKFILE" > Dockerfile
	sudo docker build . --tag $(IMG_REPO):latest
	-rm Dockerfile

# in case you also want to tag with version
ociV: build
	@echo -e "$$DOCKFILE" > Dockerfile
	sudo docker build . --tag $(IMG_REPO):$(VER) --tag $(IMG_REPO):latest
	-rm Dockerfile

clean:
	@rm -rf bin

push: oci
	sudo docker push $(IMG_REPO) $(DOCKER_REPO)

pushV: ociV
	sudo docker push $(IMG_REPO) $(DOCKER_REPO)
	sudo docker push $(IMG_REPO):$(VER) $(DOCKER_REPO):$(VER)

all: clean build oci