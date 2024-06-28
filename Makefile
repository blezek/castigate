SHELL := /bin/bash -c

build:
	go build .


test: build
	rm -rf test
	mkdir -p test
	./castigate --config test/castigate.yaml init
	./castigate --config test/castigate.yaml add 5_minutes https://5minutesinchurchhistory.ligonier.org/rss
	./castigate --config test/castigate.yaml edit --count 1 --directory foo 5_minutes
	./castigate --config test/castigate.yaml sync
