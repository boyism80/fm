-- NPC name (String.wz/Npc.img.xml): 빛나는 돌

return {
	on_click = function(me, npc)
		local q = me:quest(2166)
		if not q:started() then
			return
		end
		q:record("5")
		q:sync_progress()
		me:show_quest_completion(2166)
		me:dialog(npc, "신비로운 힘이 온 몸에 전해져 오는 것 같다.")
	end
}
