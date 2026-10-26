-- Item name (String.wz/Cash.img.xml): 눈사람과 함박눈

local weather = require("script/lib/weather")

return {
    on_cash = function(me, item_id, text)
        return weather.use(me, item_id, text)
    end,
}
