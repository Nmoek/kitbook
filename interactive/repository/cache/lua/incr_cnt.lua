-- 具体业务Key
local key = KEYS[1]
-- 哪一个数据: 阅读数、点赞数、阅读数
local cntKey = ARGV[1]

local delta = tonumber(ARGV[2])

local exist=redis.call("EXISTS", key)
if exist == 1 then

    redis.call("HINCRBY", key, cntKey, delta)
    return 1
else
    return 0
end


