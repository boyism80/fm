-- Item name (String.wz/Cash.img.xml): 해적데비존 퀘스트용

local weather = require("script/lib/weather")

return {
    on_cash = function(me, item_id, text)
        return weather.use(me, item_id, text)
    end,
}
