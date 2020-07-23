usage:
	echo "See Makefile"

check: staticcheck

install: 
	go install -v ./...

staticcheck:
	go run honnef.co/go/tools/cmd/staticcheck -- $$(go list ./... | grep -v tempfork)
