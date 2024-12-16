APP_NAME = app

.PHONY: lint
lint:
	golangci-lint run \
	-v -c .golangci.yaml --color='always' \
	--exclude-dirs-use-default --exclude-files './internal/repository/mocks/*','\*.mod','\*.sum' \
	--exclude-dirs 'vendor'

.PHONY: test
test:
	go test -v -cover ./...

#go test -v -cover ./... | grep 'ok' | sed 's/.* \([0-9.]\+\)%/\1/' | grep -oE '[0-9.]+' > coverage.txt

#go test -v -coverprofile=cover.out ./...
#go tool cover -func=cover.out -o cover.txt
#cat cover.txt | grep 'total:' | awk -F' ' '{print $$(NF)}' | sed 's/%//' > cover.txt

.PHONY: build
build:
	go build -v -o ./bin/${APP_NAME} ./cmd/${APP_NAME}

.PHONY: clean
clean:
	go clean
	rm -rf ./bin/${APP_NAME}