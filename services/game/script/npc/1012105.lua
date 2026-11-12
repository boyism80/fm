-- NPC name (String.wz/Npc.img.xml): 나미녀

local salon = require("script/lib/salon")

return {
    on_click = function(me, npc)
        salon.skin(me, npc, { coupon = 5153000 })
    end,
}
