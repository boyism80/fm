-- NPC name (String.wz/Npc.img.xml): 파풀라투스

return {
	on_click = function(me, npc)
		local npc_id = npc:id()
		local q = me:quest(6363)
		if q ~= nil and q:started() then
			if not me:dialog_yes_no(npc, "그럼 당신의 시간을 되돌릴 준비를 할게요. 이야아압~") then
				return
			end
			me:quest(6364):start(npc_id, "2")
			me:show_quest_completion(6363)
			me:map(220080000, 5)
			return
		end
		if not me:dialog_yes_no(npc, "정말 나가고 싶어요?") then
			return
		end
		me:map(220080000, 5)
	end
}
