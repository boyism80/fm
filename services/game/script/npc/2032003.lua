-- NPC name (String.wz/Npc.img.xml): 리라

local ENTRY_MAP = 211042300
local STAGE2_QUEST = 100002
local VOLCANO_BREATH = 4031062

local function ensure_quest_started(me, quest_id, record)
	local q = me:quest(quest_id)
	if q == nil then
		return nil
	end
	if q:started() or q:completed() then
		return q
	end
	if record == nil then
		record = ""
	end
	if q:wz() == nil then
		q:start(record)
	else
		q:start(0, true)
	end
	return me:quest(quest_id)
end

return {
	on_click = function(me, npc)
		if not me:dialog(npc, "저 험한 길을 어떻게 건너 온 거에요? 대단해요. #b화산의 숨결#k은 여기 있습니다. 이걸 오빠에게 전해주세요. 이제 곧 당신들이 원하는 것을 만나게 되겠군요.", true, true) then
			return
		end
		local code = me:exchange({}, {
			item = { [VOLCANO_BREATH] = 1 },
			exp = 10000,
		})
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "흐음, 인벤토리 공간이 부족하신 것 같군요.")
			return
		end
		if code ~= ExchangeResult.OK then
			return
		end
		local q = me:quest(STAGE2_QUEST)
		if q ~= nil and q:started() then
			q:force_complete(npc)
		elseif q ~= nil and not q:completed() then
			ensure_quest_started(me, STAGE2_QUEST, "")
			q = me:quest(STAGE2_QUEST)
			if q ~= nil and q:started() then
				q:force_complete(npc)
			end
		end
		me:map(ENTRY_MAP)
	end
}
