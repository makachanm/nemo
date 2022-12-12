CC = go
DATE = $(shell date +%Y%m%d%H%M)
ARCH = $(shell uname -m)

LDINFO = -X main.BuildDate=$(DATE) -X main.Arch=$(ARCH)
LDFLAGS = -s -o artifact

nemo :
	$(CC) build -ldflags "$(LDINFO)"

clean :
	rm nemo