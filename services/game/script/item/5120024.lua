-- Item name (String.wz/Cash.img.xml): 무릉도원 퀘스트용

local weather = require("script/lib/weather")

return {
    on_cash = function(me, item_id, text)
        return weather.use(me, item_id, text)
    end,
}
