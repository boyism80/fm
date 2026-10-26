-- Item name (String.wz/Cash.img.xml): 하트 확성기

local megaphone = require("script/lib/megaphone")

return {
    on_cash = function(me, item_id, text, ear)
        if me:level() < 10 then
            me:message("레벨 10 이상만 사용할 수 있습니다.", Msg.PinkText)
            return false
        end
        return megaphone.send(me, text, Msg.HeartMegaphone, MessageScope.World, ear)
    end,
}
