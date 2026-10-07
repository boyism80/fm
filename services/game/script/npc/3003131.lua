-- NPC name (String.wz/Npc.img.xml): 카오

local RECORD_QUEST = 5604004
local PETS = { 5000072, 5000073 }

return {
	on_click = function(me, npc)
		local q = me:quest(RECORD_QUEST)
		if q:started() and q:record() == "1" then
			me:dialog(npc, "이미 펫들을 데려갔는걸?")
			return
		end
		local reward = {}
		for _, pet in ipairs(PETS) do
			reward[pet] = 1
		end
		if me:exchange(nil, { item = reward }) ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족한건 아닌지 확인해보게.")
			return
		end
		if q:started() == false then
			q:start("1")
		else
			q:record("1")
		end
		me:dialog(npc, "당신 에게 신비로운 펫들을 증정하죠.")
	end
}
