doc:
	@GOPATH=$(shell pwd) godoc -http ":6060" -analysis type
