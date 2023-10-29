-- 1. 发送过来的key phone_code:[业务名称]:[手机号码]
local key = KEYS[1]

-- 2. 累计key使用次数
local cntKey = key..":cnt"
-- 3. 传入的验证码
local val = ARGV[1]

-- 4. 获取Redis中是否有该key
local ttl = tonumber(redis.call("ttl", key))
    -- key存在但是没有设置过期时间
if ttl == -1 then

    return -1
    -- key 不存在 或 key已经过了1min 则可以发送验证码
elseif ttl == -2 or ttl < 540 then

    redis.call("set", key, val)
    redis.call("expire", key, 600) --给key设置过期时间10min
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600) -- 和验证码需要一起过期
    return 0
else
    -- 发送过于频繁
    return -2
end