build: templ
    go build -o neoweb main.go

templ:
    go run github.com/a-h/templ/cmd/templ generate

check:
    go vet ./...
    golangci-lint run ./...

fmt:
    templ fmt .
    go fmt ./...

test:
    go test -v ./...


run:
    go run . -log-level debug

watch:
    go run github.com/a-h/templ/cmd/templ generate --watch --proxy="http://localhost:8080" --cmd="go run . -log-level debug"
