-- Item name (String.wz/Cash.img.xml): 해골 확성기

local megaphone = require("script/lib/megaphone")

-- The client has no skull megaphone message of its own; it shows the one-line form of the three-line megaphone.
return {
    on_cash = function(me, item_id, text, ear)
        if me:level() < 10 then
            me:message("레벨 10 이상만 사용할 수 있습니다.", Msg.PinkText)
            return false
        end
        return megaphone.send_lines(me, { text }, ear)
    end,
}
