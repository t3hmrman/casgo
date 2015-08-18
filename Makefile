all: resources install

resources:
	cd cas && rice embed-go

install:
	go install

test:
	ginkgo -r
