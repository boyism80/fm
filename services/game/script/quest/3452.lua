-- Quest name (Quest.wz/Quest.img.xml): 블록퍼스의 세계침략자?

local quest_id = 3452

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

	local file = "#fUI/UIWindow.img/QuestIcon/"
	if item_count(me, 4000099) < 1 then
		me:dialog(npc, "아직 문어 열쇠고리는 구하지 못한 모양이군. 그건 블록퍼스들이 가지고 있다네.", false, false)
		return
	end

	if item_count(me, 4001125) >= 1 then
		if not me:dialog(npc, "문어 열쇠고리는 구해 왔는가? 흐음... 귀엽게 생긴 물건이군. 하지만 이게 바로 지구를 위협하는 외계인의 정체를 밝힐 중요한 물건...잠깐!", false, true) then
			return
		end
		if not me:dialog(npc, "자네가 가지고 있는 그 물건! 그것 좀 보여주게. 자네 손에 들고 있는 바로 그 설계도 말일세. 오~ 이런 물건을 어디서 구한건가? 이것만 있으면 블록퍼스에 대한 연구를 더 빨리 진행시킬 수 있겠어.", false, true) then
			return
		end
		if not me:dialog(npc, "뜻밖의 수확인걸. 좋아. 자네에게 특별한 선물을 하도록 하지. 분명 도움이 될거야. 하하하~\r\n\r\n" .. file .. "4/0#\r\n#v2040701# #t2040701# 1개\r\n\r\n" .. file .. "8/0# 16000 exp", false, true) then
			return
		end
		local code = me:exchange(
			{ item = { [4000099] = 1, [4001125] = 1 } },
			{ item = { [2040701] = 1 }, exp = 16000 }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "뭘 그렇게 많이 들고 다니는건가? 인벤토리에 빈 칸이 있는지 확인해 주게.", false, false)
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "블록퍼스와 외계인... 아무리 생각해도 뭔가 비슷하단 말이야. 아니, 문어와 외계인이 비슷한 걸지도... 흐음. 이것도 새로운 이론이군.", false, false)
		q:force_complete(npc)
	else
		if not me:dialog(npc, "문어 열쇠고리는 구해 왔는가? 흐음... 귀엽게 생긴 물건이군. 하지만 이게 바로 지구를 위협하는 외계인의 정체를 밝힐 중요한 물건이지. 정말 고맙네.\r\n\r\n" .. file .. "4/0#\r\n#v2000011# #t2000011# 50개\r\n\r\n" .. file .. "8/0# 8000 exp", false, true) then
			return
		end
		local code = me:exchange(
			{ item = { [4000099] = 1 } },
			{ item = { [2000011] = 50 }, exp = 8000 }
		)
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "뭘 그렇게 많이 들고 다니는건가? 인벤토리에 빈 칸이 있는지 확인해 주게.", false, false)
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "블록퍼스와 외계인... 아무리 생각해도 뭔가 비슷하단 말이야. 아니, 문어와 외계인이 비슷한 걸지도... 흐음. 이것도 새로운 이론이군.", false, false)
		q:force_complete(npc)
	end
end
