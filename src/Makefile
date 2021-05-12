run:
	@go run .

build:
	@go build -ldflags "-s -w"

buildlinux:
	@mkdir -p dist

	@GOOS=linux GOARCH=amd64 go build -o inventory -ldflags "-s -w"
	@tar cvzf dist/inventory-linux-amd64.tar.gz public/ templates/ inventory ../LICENSE ../README.md
	@if [ -f "inventory" ]; then rm inventory; fi

	@GOOS=linux GOARCH=arm64 go build -o inventory -ldflags "-s -w"
	@tar cvzf dist/inventory-linux-arm64.tar.gz public/ templates/ inventory ../LICENSE ../README.md
	@if [ -f "inventory" ]; then rm inventory; fi

clean:
	@if [ -d "dist" ]; then rm -r dist; fi
	@if [ -f "inventory" ]; then rm inventory; fi
	@echo "cleaning"

