-- 스탠의 두 번째 가르침 (Quest.wz/QuestData/6033.img.xml): 스탠의 두 번째 가르침

local quest_id = 6033

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() then
		q:start(npc, true)
		return
	end

	if not me:dialog_yes_no(npc, "그래 가져왔군. 어디 한번 볼까?") then
		me:dialog(npc, "뭐? 보여주기 싫다고? 그럼 나도 자네를 도와줄 이유는 없지.", false, false)
		return
	end

	local data = q:record()
	if data == nil or data == "" then
		if not me:dialog(npc, "엉터리군! 꽝이야! 이런 식으로 물건을 만드니까 문제가 생기는거다! 다시 만들어와!", false, true) then
			return
		end
		local code = me:exchange({ item = { [4260003] = 1 } }, {})
		if code ~= ExchangeResult.OK then
			return
		end
		q:record("1")
		return
	end

	if data == "1" then
		local file = "#fUI/UIWindow.img/QuestIcon/"
		if not me:dialog(npc, "뭐. 썩 맘에 들지는 않지만 일단은 통과다!! 앞으로도 절대 해이해지기 않도록 주의하라고!!\r\n\r\n" .. file .. "4/0#\r\n\r\n#fSkill/000.img/skill/0001007/icon# #q1007# (레벨 2)\r\n\r\n" .. file .. "8/0# 230000 exp", false, true) then
			return
		end
		local skill = me:add_skill(1007)
		if skill ~= nil then
			skill:level(2)
		end
		local code = me:exchange({ item = { [4260003] = 1 } }, { exp = 230000 })
		if code ~= ExchangeResult.OK then
			return
		end
		q:force_complete(npc)
		me:show_effect(EffectType.QuestCompletion)
		return
	end

	me:dialog(npc, data, false, false)
end
