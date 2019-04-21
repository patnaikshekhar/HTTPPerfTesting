$env:GOOS = "windows"
go build
.\perftest.exe --concurrency 5 --number 20 http://localhost:8080