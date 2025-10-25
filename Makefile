lint:
	revive -config revive.toml -set_exit_status ./...

test:
	go test -v ./...

qa: lint test

install:
	go install ./cmd/zettelkasten
