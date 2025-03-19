install-devtools:
	go install github.com/g4s8/envdoc@v1.4.0
	go install github.com/a-h/templ/cmd/templ@latest

generate-templ:
	templ generate

generate-onrun: generate-templ

generate: generate-templ
	go generate ./...
