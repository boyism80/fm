-- NPC name (String.wz/Npc.img.xml): 사자왕의 수호석상

local ex = require("script/lib/expedition")

local EXIT_MAP = 211061001
local BATTLE_MAPS = {
	[211070100] = true,
	[211070101] = true,
	[211070110] = true,
}

return {
	on_click = function(me, npc)
		if BATTLE_MAPS[me:map():wz():id()] then
			if me:dialog_yes_no(npc, "여기서 나가시겠습니까?") then
				me:map(EXIT_MAP)
			end
			return
		end
		ex.talk(me, npc, "von_leon")
	end
}
