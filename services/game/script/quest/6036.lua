-- 뜻밖의 결과 (Quest.wz/QuestData/6036.img.xml): 뜻밖의 결과

local quest_id = 6036

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() then
		q:start(npc, true)
		return
	end

	local file = "#fUI/UIWindow.img/QuestIcon/"
	if not me:dialog_yes_no(npc, "엇! 자네! 이걸 정말 자네가 만들었단 말인가! 이건 내가 그렇게 만들고 싶어하던 황금모루!! 이걸 나에게 주지 않겠나?\r\n\r\n" .. file .. "4/0#\r\n\r\n#fSkill/000.img/skill/0001007/icon# #q1007# (레벨 3)\r\n\r\n" .. file .. "8/0# 960000 exp") then
		return
	end

	me:dialog(npc, "이거 정말 굉장하군 굉장해! 고맙네. 고마워. 이런 물건까지 만들어내다니 이제 더 이상 내가 가르칠 게 없겠군. 자네를 제작 마스터로 인정하지! 하하!!", false, false)
	local skill = me:add_skill(1007)
	if skill ~= nil then
		skill:level(3)
	end
	local code = me:exchange({ item = { [4031980] = 1 } }, { exp = 960000 })
	if code ~= ExchangeResult.OK then
		return
	end
	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
end
