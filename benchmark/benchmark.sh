wrk -t8 -c1000 -d15s -s test.lua http://127.0.0.1:8080/api/v2/agent/log >report.txt
