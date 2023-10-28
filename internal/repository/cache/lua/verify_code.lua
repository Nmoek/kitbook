local key = KEY[1]
local cntKey = key.."cnt"

-- 用户输入的验证码
local expectedCode = ARGV[1]

--当前还剩下几次验证机会
local curCnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)

-- 验证码次数字段莫名其妙失效了
if curCnt == nil then
    return -1
-- 超出验证次数
elseif curCnt <= 0 then
    return -2
end

if code == expectedCode then
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    return -3
end