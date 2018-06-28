lint:
	gometalinter .
	gometalinter controller/.
	gometalinter entity/.
	gometalinter errors/.
	gometalinter handlers/. --disable gocyclo
	gometalinter postgres/.

build:
	go build -o bin/game main.go

test:
	go test github.com/Tournament/handlers/.
	go test github.com/Tournament/postgres/.

run:
	bin/game

dockerbuild:
	docker build -t tournament .

dockerrun:
	docker run --rm --name tournament -p 8080:8080 --net=host tournament
	