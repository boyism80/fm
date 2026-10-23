-- Item name (String.wz/Cash.img.xml): 세줄 확성기

local megaphone = require("script/lib/megaphone")

return {
    on_cash = function(me, item_id, text, ear, pet_sn, input)
        return megaphone.send_lines(me, input.lines, ear)
    end,
}
