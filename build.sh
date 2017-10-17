#CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' ./cmd/facter/
CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' ./cmd/facter/
