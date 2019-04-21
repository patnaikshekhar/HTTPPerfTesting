$env:GOOS = "linux"
$env:CGO_ENABLED = 0

go build
docker build -t patnaikshekhar/perftest-go-server:1 .