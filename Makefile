.PHONY: build

build: build/ascanvas

clean:
	rm build/ascanvas

build/ascanvas:
	commit="dev"
	date=`date +%FT%T%z`
	go build -ldflags "-w -extldflags '-static' -X=main.appCommit=$commit -X=main.appBuilt=$date" -o build/ascanvas github.com/fluxynet/ascanvas/cmd/ascanvas