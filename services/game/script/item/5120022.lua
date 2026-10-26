-- Item name (String.wz/Cash.img.xml): 줄리엣 퀘스트용

local weather = require("script/lib/weather")

return {
    on_cash = function(me, item_id, text)
        return weather.use(me, item_id, text)
    end,
}
