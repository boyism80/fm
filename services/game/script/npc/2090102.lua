-- NPC name (String.wz/Npc.img.xml): 나란

local salon = require("script/lib/salon")

return {
    on_click = function(me, npc)
        salon.skin(me, npc, { coupon = 5153006 })
    end,
}
