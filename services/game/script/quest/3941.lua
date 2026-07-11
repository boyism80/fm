-- Quest name (Quest.wz/Quest.img.xml): 왕비의 비단 빼앗기

local quest_id = 3941

local function morph_source(me)
	local b = me:buff(BuffFlag.Morph)
	if not b then
		return -1
	end
	local wz = b:wz()
	if type(wz.id) == "number" then
		return wz.id
	end
	return wz:id()
end

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if morph_source(me) ~= 2210005 then
		me:dialog(npc, "날씨가 정말 덥죠? 하지만 또 더워야 진짜 사막이라고 할 수 있지 않겠어요? 그나저나 #p2101004#님은 언제 오시려나?", false, false)
		return
	end

	local sel = me:dialog_list(npc, "...#p2101004#님 아니십니까? 오랜만입니다. 이번에 왕비님께서 애타게 찾으시던 비단을 다행히 구해놨습니다. 항상 그렇듯이 물건은 최고급... 그런데 왜 이렇게 땀을 흘리십니까?\r\n", {
		"(음성변조)아니, 그냥 태양이 더워서...",
	})
	if sel == nil then
		return
	end

	sel = me:dialog_list(npc, "아리안트가 언제는 안 더운 곳이었습니까? 항상 이랬지만 #p2101004#님은 더위를 안 타시는 줄 알았습니다만... 아니 점점 얼굴이 빨개지고 계십니다. 괜찮으십니까?\r\n", {
		"(음성변조)괘, 괜찮네. 걱정하지 말게...",
	})
	if sel == nil then
		return
	end

	sel = me:dialog_list(npc, "정말 괜찮으신 겁니까? 안 그래도 #p2101004#님은 허약체질 같으신데, 역시 약이라도 좀 드셔야 하는 것 아닙니까? 엘나스 쪽에서 구한, 감기약이 있는데 하나 사시겠습니까? 싸게 드리죠.\r\n", {
		"괘, 괜찮다니까 그러네!",
	})
	if sel == nil then
		return
	end

	if not me:dialog_accept(npc, "괜찮으시다고요? 그런데 #p2101004#님 목소리가 평소하고 아주 많이 다릅니다? 정말 감기에 걸리신 것 아닙니까? 도무지 평소의 #p2101004#님 같지 않으십니다. 평소에는 항상 #t4010007#을 싸게 달라고 하셨는데... 이상하군요. 정말 #p2101004#님이십니까?") then
		me:dialog(npc, "#p2101004#님이 아니라고요? 아니 어딜 봐도 #p2101004#님이신데 도대체 무슨 말을 하시는 건지... 정말 많이 아프신 모양이군요. 비단은 괜찮아 지면 드릴 테니, 다시 찾아오시죠.", false, false)
		return
	end
	me:dialog(npc, "평소의 #p2101004#님 같지 않으십니다. 원래 #p2101004#님이 이렇게 말이 적은 분이 아니신 걸로 아는데... 무슨 일이라도 있으신 겁니까? 헙. 얼굴이 점점 붉어지시는 게 아무래도 화가 나신 모양이군요. 죄송합니다. 어서 가서 비단을 가져오겠습니다. 잠시만 기다려 주세요.", false, false)
	q:start(npc, true)
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

	if morph_source(me) ~= 2210005 then
		me:dialog(npc, "날씨가 정말 덥지 않아요? 목이 바짝바짝 타는 군요.", false, false)
		return
	end

	local file = "#fUI/UIWindow.img/QuestIcon/"
	if not me:dialog(npc, "자 여기 있습니다. 조심해서 가져가십시오. 이 비단, 구하기 무척 힘든 물건입니다. 혹시 조금이라도 찢어진 곳이 있으면 당장 왕비님께서 #p2101004#님을 감옥에 가두고 말 것입니다.\r\n\r\n" .. file .. "4/0#\r\n#v4031571# #t4031571# 1개\r\n\r\n", false, true) then
		return
	end

	local code = me:exchange({}, { item = { [4031571] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리에 짐이 많아서 비단을 넣어드릴 수가 없는걸요? 중요한 물건을 운반하시는데 조금만 자리를 비워보세요.", false, false)
		return
	end
	if code ~= ExchangeResult.OK then
		return
	end
	q:force_complete(npc)
end
