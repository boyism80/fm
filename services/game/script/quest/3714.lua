-- Quest name (Quest.wz/Quest.img.xml): 혼테이즈에게 남긴 돌...

local quest_id = 3714

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog(npc, "뀨웃....뀨르르르르...\r\n\r\n\#b(앗 귀엽다! 가까이 다가가 보았다.)#k", false, true) then
			return
		end
		local code = me:exchange({ item = { [4001094] = 1 } }, { item = { [2041200] = 1 } })
		if code == ExchangeResult.LackCapacity then
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "아기용이 붉은색 돌을 뱉어냈습니다.", false, true)
		q:start(npc, true)
		q:force_complete(npc)
	end
}
