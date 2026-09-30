-- NPC name (String.wz/Npc.img.xml): 무스

return {
	on_click = function(me, npc)
		local q = me:quest(6180)
		if q == nil or not q:started() then
			me:dialog(npc, "방패 수련장? 어디서 그걸 들은거야?", false, true)
			return
		end
		if not me:dialog(npc, "좋아. 방패 수련장으로 보내주도록 하지.", false, true) then
			return
		end
		me:map(924000000, 0)
	end
}
