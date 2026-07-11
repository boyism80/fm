local quest_id = 2186

local function item_count(me, item_id)
	local slots = me:item(item_id)
	local count = 0
	for _, it in pairs(slots) do
		count = count + it:count()
	end
	return count
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not q:started() then
		q:start(npc, true)
		return
	end

	local c1853 = item_count(me, 4031853)
	if c1853 < 1 then
		return
	end

	local have_diff = 0
	local c1854 = item_count(me, 4031854)
	local c1855 = item_count(me, 4031855)
	if c1854 >= 1 then
		have_diff = have_diff + 1
	end
	if c1855 >= 1 then
		have_diff = have_diff + 1
	end

	local str = "앗! 제 안경을 찾으셨다고요? 한번 써볼까요? 제 안경이 아닐 수도 있으니 직접 써보는 방법 말고는 구분할 수 없어요."
	local reward_count = 5
	local reward_exp = 1000
	if have_diff == 0 then
		str = str .. "제 안경이 맞군요. 감사합니다.\r\n\r\n"
		str = str .. "#fUI/UIWindow.img/QuestIcon/4/0#\r\n"
		str = str .. "#i2030019# #t2030019# 5개\r\n\r\n"
		str = str .. "#fUI/UIWindow.img/QuestIcon/8/0# 1000 exp\r\n"
		reward_count = 5
		reward_exp = 1000
	elseif have_diff == 1 then
		str = str .. "어라? 손에 들고 계신 또 다른 안경은 무엇인가요? 어? 그 안경도 제게 주시지 않으시겠어요? 보상은 해드릴게요. 후후, 기분에 따라서 안경을 바꿔서 써보기도 해야겠어요.\r\n\r\n"
		str = str .. "#fUI/UIWindow.img/QuestIcon/4/0#\r\n"
		str = str .. "#i2030019# #t2030019# 15개\r\n\r\n"
		str = str .. "#fUI/UIWindow.img/QuestIcon/8/0# 1500 exp\r\n"
		reward_count = 15
		reward_exp = 2000
	else
		str = str .. "어라? 손에 들고 계신 또 다른 안경들은 무엇인가요? 어? 그 안경들도 제게 주시지 않으시겠어요? 보상은 해드릴게요. 후후, 기분에 따라서 안경들을 바꿔서 써보기도 해야겠어요.\r\n\r\n"
		str = str .. "#fUI/UIWindow.img/QuestIcon/4/0#\r\n"
		str = str .. "#i2030019# #t2030019# 30개\r\n\r\n"
		str = str .. "#fUI/UIWindow.img/QuestIcon/8/0# 2000 exp\r\n"
		reward_count = 30
		reward_exp = 2000
	end

	me:dialog(npc, str, false, true)
	me:dialog(npc, "후후훗... 다시 낚시를 즐겨 볼까요?", false, true)

	local cost_items = { [4031853] = c1853 }
	if c1854 > 0 then
		cost_items[4031854] = c1854
	end
	if c1855 > 0 then
		cost_items[4031855] = c1855
	end
	local code = me:exchange(
		{ item = cost_items },
		{ item = { [2030019] = reward_count }, exp = reward_exp }
	)
	if code == ExchangeResult.LackCapacity then
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end

	q:force_complete(npc)
	me:show_effect(EffectType.QuestCompletion)
	local map = me:map()
	if map ~= nil then
		for _, n in pairs(map:npcs()) do
			if n:id() == npc then
				n:show_effect("quest")
				break
			end
		end
	end
end
