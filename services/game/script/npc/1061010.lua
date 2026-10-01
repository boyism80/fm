-- NPC name (String.wz/Npc.img.xml): 빛나는 수정

local exits = {
	[108010301] = 105070001,
	[108010201] = 100040106,
	[108010101] = 105040305,
	[108010401] = 107000402,
	[108010501] = 105070200,
}

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "정말 이곳에서 나가시겠습니까?") then
			return
		end
		local map = me:map()
		if map == nil then
			return
		end
		local dest = exits[map:template_id()]
		if dest == nil then
			return
		end
		me:map(dest)
	end
}
