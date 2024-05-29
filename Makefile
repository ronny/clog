all: test

test:
	go test -shuffle=on -coverprofile=coverage.txt -count=1 ./...

bench:
	go test -run=XXX -benchmem -count=1 -bench=. ./...

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

format:
	gofumpt -l -w .

lint: vet staticcheck format

install-tools:
	go install mvdan.cc/gofumpt@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

ci-checks: lint
	go mod tidy
	git diff --exit-code

ci-test: install-tools ci-checks test

# Tell Make that these steps should always run, never cached, even if we have
# files with the same names as these in the filesystem.
.PHONY: all test vet staticcheck lint format mod-tidy ci-checks install-tools ci-test
