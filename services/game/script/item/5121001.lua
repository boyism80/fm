-- Item name (String.wz/Cash.img.xml): 필승코리아 응원하기

local weather = require("script/lib/weather")

return {
    on_cash = function(me, item_id, text)
        return weather.use(me, item_id, text)
    end,
}
