-- Quest name (custom Resonance): [암벽 거인] 수상한 움직임

local quest_id = 31348

return {
	on_start = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog_yes_no(npc, "내가 시력이 좋은 건 알고 있지? 암벽 거인의 머리 위쪽에서 무언가 수상한 사람들의 움직임을 발견했어. 암벽 거인의 말에 따르면... 그들이 암벽 거인을 오염 시키려는 주범들이 아닐까? 지금 너가 그쪽으로 올라가서 수상한 사람들이 없는지 살펴줘") then
			me:dialog(npc, "정말로 머쉬킹 왕국에 도움을 줘볼 생각이 없나? 언제든지 생각이 바뀌면 나를 찾아오게나.", false, false)
			return
		end

		local from = me:map()
		local from_id = from:wz().id
		me:map(924030000)
		local m = me:map()
		if m == nil then
			return
		end
		local others = 0
		for _, ch in pairs(m:characters()) do
			if ch:id() ~= me:id() then
				others = others + 1
			end
		end
		if others > 0 then
			me:notice("누군가 안에 퀘스트를 진행중입니다", Msg.PinkText)
			me:map(from_id)
			return
		end

		m:kill_all_mobs()
		for _ = 1, 6 do
			m:spawn_mob(9100044, 674, 60)
		end
		q:start(npc, true)
		me:open_npc(2210009)
	end,

	on_end = function(me, npc)
		local q = me:quest(quest_id)
		if q == nil then
			return
		end

		if not me:dialog(npc, "암벽거인은 무사할까...?", false, true) then
			return
		end
		if not me:dialog_yes_no(npc, "그들의 말대로라면 분명히 지금쯤 위기에 처해있을 터... 암벽 거인에게로 가서 대화해보자.") then
			return
		end
		q:force_complete(npc)
		local next_q = me:quest(31349)
		if next_q ~= nil then
			next_q:start(npc, true)
		end
	end
}
