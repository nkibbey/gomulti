LEN := 0

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

build:
	@mkdir -p bin
	go build -ldflags="-X 'main.Version=$(VER)' -X 'main.GitCommit=$(CI_COMMIT_SHA)' -X 'main.BuildTime=$$(date)'" -tags netgo -o bin ./...

clean:
	@rm -rf bin

all: clean build