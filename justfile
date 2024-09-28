clideps:
    go install github.com/a-h/templ/cmd/templ@latest

build: templ
    go build -o neoweb main.go

templ: clideps
    templ generate

check:
    go vet ./...
    golangci-lint run ./...

fmt:
    templ fmt .
    go fmt ./...

test:
    go test -v ./...


watch:
    templ generate --watch --proxy="http://localhost:8080" --cmd="go run . -log-level debug"
