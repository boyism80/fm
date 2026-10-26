-- Item name (String.wz/Cash.img.xml): 생일 축하해

local weather = require("script/lib/weather")

return {
    on_cash = function(me, item_id, text)
        return weather.use(me, item_id, text)
    end,
}
