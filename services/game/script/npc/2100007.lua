-- NPC name (String.wz/Npc.img.xml): 라일라

local salon = require("script/lib/salon")

return {
    on_click = function(me, npc)
        salon.skin(me, npc, 5153007)
    end,
}
