CC = go
DATE = $(shell date +%Y%m%d%H%M)
ARCH = $(shell uname -m)

LDINFO = -X main.BuildDate=$(DATE) -X main.Arch=$(ARCH)
LDFLAGS = -s -o artifact
OPFLAGS = -trimpath 

nemo :
	$(CC) build $(OPFLAGS) -ldflags "$(LDINFO)"

clean :
	rm nemo