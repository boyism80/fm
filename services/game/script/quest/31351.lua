-- [암벽 거인] 암벽 거인 구출 작전 7 (Quest.wz/QuestInfo.img.xml): [암벽 거인] 암벽 거인 구출 작전 7

local quest_id = 31351

function on_start(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog_yes_no(npc, "찐득하고 소름 끼치는 무언가가 조금씩 나를 갉아먹는다. 이대로라면 결국 나는 그들에게 조종당하고 만다.\r\n마지막 부탁이다. 나의 몸 속에서도 가장 깊은 곳... #b심장#k에 들어가 무엇이 있나 확인해다오. 분명 그 곳에 무언가가 있을 거다. 서둘러다오.") then
		return
	end

	q:start(npc, true)
end

function on_end(me, npc)
	local q = me:quest(quest_id)
	if q == nil then
		return
	end

	if not me:dialog(npc, "타란튤로스를 해치운건가?", false, true) then
		return
	end
	if not me:dialog_yes_no(npc, "타란튤로스를 해치우자 암벽 거인의 떨림이 잦아들고 맑은 기운이 사방에서 샘솟기 시작한다. 암벽 거인들 위기로부터 구출해냈다. 이제 당분간은 오염되지 않을 것 같다.") then
		return
	end

	me:map(240092100)
	q:force_complete(npc)
end
