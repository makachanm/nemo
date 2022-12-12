VERSION = 0.9.0

CC = go
DATE = $(shell date +%Y%m%d%H%M)
ARCH = $(shell uname -m)

LDINFO = -X main.BuildDate=$(DATE) -X main.Arch=$(ARCH) -X main.Version=$(VERSION)
LDFLAGS = -s -o artifact
OPFLAGS = -trimpath 

nemo :
	$(CC) build $(OPFLAGS) -ldflags "$(LDINFO)"

clean :
	rm nemo