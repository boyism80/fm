-- NPC name (String.wz/Npc.img.xml): 지나

local salon = require("script/lib/salon")

return {
    on_click = function(me, npc)
        salon.skin(me, npc, 5153002)
    end,
}
