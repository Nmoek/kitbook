--限流对象(IP组成的key)
local key = KEYS[1]
--窗口大小
local windowSize = tonumber(ARGV[1])
--限流阈值
local threshold = tonumber(ARGV[2])
--当前时间戳
local now = tonumber(ARGV[3])
--窗口下限
local min = now - windowSize

--删除不在窗口范围内的数据
redis.call("ZREMRANGEBYSCORE", key, "-inf", min)
--计算当前窗口中请求数
--local cnt = redis.call("ZCOUNT", key, min, now) --这种写法可能存在时间卡不准的情况
local cnt = redis.call("ZCOUNT", key, "-inf", "+inf")
if cnt >= threshold then

--    需要限流
    return 'true'
else
--    不用限流
    redis.call("ZADD", key, now, now)
    redis.call("PEXPIRE", key, windowSize)
    return 'false'
end