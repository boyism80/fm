-- NPC name (String.wz/Npc.img.xml): 줄리엣

local rj = require("script/lib/romeo_juliet")

return {
	on_click = function(me, npc)
		rj.strip_items(me)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		local map_id = map:wz().id
		if map_id >= 926100000 and map_id < 926110000 then
			me:map(926100700)
		else
			me:map(926110700)
		end
	end
}
