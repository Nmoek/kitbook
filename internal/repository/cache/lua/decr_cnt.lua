local key = KEYS[1]
local cntKey = ARGS[1]
local delta =  tonumber(ARGS[2])

local exists=redis.call("EXISTS", key)
if exists == 1 then

    redis.call("HINCRBY", key, cntKey,-delta)
    return 1
else
    return 0
end