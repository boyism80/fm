-- NPC name (String.wz/Npc.img.xml): 과일가게 할아버지

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		local map_id = wz:id()
		if not me:dialog_yes_no(npc, "정말 이곳에서 나가 #b#m" .. (map_id - 1) .. "##k 맵으로 돌아가고 싶나?") then
			return
		end
		me:map(map_id - 1)
	end
}
