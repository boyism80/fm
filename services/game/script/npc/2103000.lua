-- NPC name (String.wz/Npc.img.xml): 궁전 오아시스

return {
	on_click = function(me, npc)
		local q = me:quest(3900)
		if q ~= nil and q:started() and q:record() ~= "5" then
			me:exchange(nil, { exp = 300 })
			q:record("5")
			q:sync_progress()
			me:show_quest_completion(3900)
			me:notice("오아시스의 물을 마셨습니다.", Msg.PinkText)
		end
	end
}
