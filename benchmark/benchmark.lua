-- Generate a random number and convert to string
local randomNumber = math.random(1, 100000000)
local randomString = tostring(randomNumber)

wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] =
	"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoidGVzdCIsImlzX3VzZXJfdG9rZW4iOmZhbHNlfQ.DIDe7iDi8QkakGD6aQHtiI_729Y_rbgPMnYgbWSpBac"
wrk.body = string.format('{"test":"test","something else":{"nested":1}, "unique": %s}', randomString)
