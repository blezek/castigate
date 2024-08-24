SHELL := /bin/bash -c

build:
	go test ./...
	go build .

release:
	goreleaser release --clean

release-test:
	goreleaser release --skip=publish --clean --snapshot


test: build
	rm -rf test
	mkdir -p test
	./castigate --config test/castigate.yaml init
	./castigate --config test/castigate.yaml add 5_minutes https://5minutesinchurchhistory.ligonier.org/rss
	./castigate --config test/castigate.yaml add --count 1 boys http://minecraft.blezek.com:3333/rss/6953
	./castigate --config test/castigate.yaml edit --count 1 --directory foo 5_minutes
	./castigate --config test/castigate.yaml sync
