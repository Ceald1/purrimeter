wrk -t10 -c5000 -d60s -s test.lua http://127.0.0.1:8080/api/v2/agent/log >report.txt
