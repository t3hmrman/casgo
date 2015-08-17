resources:
	go-bindata -o _resources.go public/... templates/...

all: resources
