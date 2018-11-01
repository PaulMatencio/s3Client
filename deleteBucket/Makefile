# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
#BINARY_NAME=putObject 
BINARY_UNIX=$(BINARY_NAME)_unix
      
install: 
	@echo "+ Installing $@  $(BINARY_NAME)"
	$(GOINSTALL)  
build: 
	@echo "+ Building $@  $(BINARY_NAME)"
	$(GOBUILD)   
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	@echo "+ Removing $@  $(BINARY_NAME)"
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
deps:
	$(GOGET) -u github.com/minio/minio-go 
	$(GOGET) -u github.com/s3Client/lib 
    
