build: templ
    go build -o neoweb main.go

templ:
    templ generate

check:
    go vet ./...

fmt:
    templ fmt .
    go fmt ./...

test:
    go test -v ./...


watch:
    templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."
