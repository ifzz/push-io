--
-- Created by IntelliJ IDEA.
-- User: zhibinpan
-- Date: 14/2/2017
-- Time: 5:57 PM
-- To change this template use File | Settings | File Templates.
--

request = function()
    path = "/api/v1/notification"
    wrk.headers["Content-Type"] = "application/json; charset=utf-8"
    wrk.body = "{\"appId\":\"gftrader\",\"appKey\":\"1163CFFD87155CD634CBD3DA9F53D\",\"topic\": \"mike\",\"message\":{\"msg\":{\"top\":\"123\"}, \"type\":\"quoteStock\"}}"
    return wrk.format("POST", path)
end
