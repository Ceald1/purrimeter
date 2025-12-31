-- Make counter thread-local to avoid race conditions
local counter = 0
local thread_id = 0

function setup(thread)
	-- Give each thread a unique ID
	thread_id = thread.id or math.random(1, 10000)
end

request = function()
	counter = counter + 1

	-- Create truly unique ID combining thread and counter
	local unique_id = string.format("%d-%d-%d", thread_id, counter, os.time() * 1000 + counter % 1000)

	-- Vary the log content to ensure different hashes
	local body = string.format(
		'{"log_id": "%s","timestamp": %d,"test": "test_%d","something_else": {"nested": %d},"unique": "%s","random": %f}',
		unique_id,
		os.time() * 1000 + counter,
		counter,
		counter % 100,
		unique_id,
		math.random()
	)

	return wrk.format("POST", "/api/v2/agent/log", {
		["Content-Type"] = "application/json",
		["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImlzX3VzZXJfdG9rZW4iOmZhbHNlfQ.DIDe7iDi8QkakGD6aQHtiI_729Y_rbgPMnYgbWSpBac",
	}, body)
end

response = function(status, headers, body)
	if status ~= 200 then
		print("Response:", status, body)
	end
end
