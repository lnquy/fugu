APP_NAME = "fugu"
PACKAGES = $(shell go list ./... | grep -v /vendor/)

submit: all format

all: test vet build

run: clean build_fe run_be

vet:
	go vet $(PACKAGES)

format:
	go fmt $(PACKAGES)

test:
	go test -cover -race -v $(PACKAGES)

bench:
	go test -bench=. $(PACKAGES)

build:
	go build -o $(APP_NAME) main.go

build_fe:
	rm -rf templates/*; \
	cd frontend/; \
	npm run build; \
	mv dist/index.html ../templates/index.html; \
	mkdir -p ../templates/static/; \
	mv dist/static/* ../templates/static/; \
	rm -rf dist; \
	cd ..

run_be:
	make build
	./$(APP_NAME)

clean:
	rm -f go-playground*; \
	rm -rf templates/*; \
	rm -f frontend/index_*.html; \
	rm -rf frontend/dist

.PHONY:
	submit \
	all \
	run \
	vet \
	format \
	test \
	bench \
	build \
	build_fe \
	run_be \
	clean
