wrk -t8 -c1000 -d15s -s benchmark.lua http://localhost:8080/api/v2/agent/log >report.txt
