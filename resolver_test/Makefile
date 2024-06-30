CLIENTNAME=client
SERVERNAME=server
CLIENTFILE=./client
SERVERFILE=./server
BINDIR=bin
GOBUILD=CGO_ENABLED=0 go build -gcflags=all='-N -l'
BUFGEN=buf generate

all: bufgen resolver_client resolver_server

bufgen:
	$(BUFGEN)

resolver_client:
	$(GOBUILD) -o $(BINDIR)/$(CLIENTNAME) $(CLIENTFILE)

resolver_server:
	$(GOBUILD) -o $(BINDIR)/$(SERVERNAME) $(SERVERFILE)

clean:
	rm $(BINDIR)/*