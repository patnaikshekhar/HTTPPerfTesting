docker rm -vf perftest
docker run --name perftest -p 8080:8080 patnaikshekhar/perftest-go-server:1