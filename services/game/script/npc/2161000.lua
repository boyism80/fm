-- NPC name (String.wz/Npc.img.xml): 반 레온

return {
	on_click = function(me, npc)
		if not me:dialog_accept(npc, "나를 물리치러온 모험가인가.. 하긴.. 너가 누구던 상관은 없다. 이미 서로에게 목적이 있다면 왈가 불가할 이유는 없겠지 자 덤벼라.") then
			return
		end
		local map = me:map()
		if next(map:mobs()) ~= nil then
			return
		end
		map:remove_npc(npc)
		map:spawn_mob(8840010, 0, -181)
	end
}
