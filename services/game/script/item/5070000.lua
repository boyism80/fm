-- Item name (String.wz/Cash.img.xml): 일반 확성기

local megaphone = require("script/lib/megaphone")

return {
    on_cash = function(me, item_id, text)
        return megaphone.send(me, text, Msg.Megaphone, MessageScope.Map)
    end,
}
