--
-- Created by IntelliJ IDEA.
-- User: zhibinpan
-- Date: 14/2/2017
-- Time: 5:40 PM
-- To change this template use File | Settings | File Templates.
--

request = function()
    path = "/api/v1/notification"
    wrk.headers["Content-Type"] = "application/json; charset=utf-8"
    wrk.body = "{\"appId\":\"demo\",\"appKey\":\"demo\",\"topic\": \"test\",\"message\":{\"msg\":{\"top\":\"123\"}, \"type\":\"quoteStock\"}}"
    return wrk.format("POST", path)
end
