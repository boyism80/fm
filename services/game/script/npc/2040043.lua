-- NPC name (String.wz/Npc.img.xml): 블루 벌룬

local pq = require("script/lib/party_quest")

local REPEAT_QUEST = 199600

local function map_stage(map_id)
	if map_id == nil then
		return 0
	end
	return math.floor((map_id % 922010000) / 100)
end

local function ensure_stage(sm)
	if sm:get_property("stage") == "" then
		sm:set_property("stage", "1")
	end
	return tonumber(sm:get_property("stage")) or 1
end

local function give_exp(sm, amount)
	if sm == nil or amount == nil or amount <= 0 then
		return
	end
	for _, p in ipairs(sm:players()) do
		if p ~= nil then
			local give = amount
			local q = p:quest(REPEAT_QUEST)
			local count = 0
			if q ~= nil then
				count = tonumber(q:record()) or 0
			end
			if count > 0 then
				give = math.floor(amount * 70 / 100)
			end
			p:exchange({}, { exp = give })
		end
	end
end

local function area_pattern(map)
	local pos = ""
	local total = 0
	for i = 0, 8 do
		local n = map:players_in_area(i)
		total = total + n
		pos = pos .. tostring(n)
	end
	return pos, total
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:dialog(npc, "오류가 발생했어요.")
			return
		end
		local map = me:map()
		if map == nil then
			return
		end
		local stage = ensure_stage(sm)
		local wz = map:wz()
		if wz == nil then
			return
		end
		local cur = map_stage(wz.id)
		if stage > cur then
			me:dialog(npc, "포탈이 열렸어요~ 다음 스테이지로 이동해 주세요!")
			return
		end
		local guide = sm:get_property("guideRead")
		if guide ~= "s" or not pq.is_leader(me) then
			if pq.is_leader(me) then
				sm:set_property("guideRead", "s")
				sm:set_property("stage8rand", pq.shuffle("000011111"))
			end
			me:dialog(npc, "여덟번째 스테이지에 대해 설명해 드리겠습니다. 이곳에는 여러 개의 발판이 있습니다. 이 발판 중에서 #b5개가 다음 스테이지로 향하는 포탈#k과 통해 있습니다. 파티원 중에서 #b5명이 정답 발판을 찾아 위에 올라서면#k 됩니다.\r\n단, 발판 끝에 아슬아슬하게 걸쳐서 서지 말고 발판 중간에 서야 정답으로 인정되니 이점 주의해 주시기 바랍니다. 그리고 반드시 5명만 발판 위에 올라가 있어야 합니다. 파티원이 발판에 올라서면 파티장은 #b저를 더블클릭하여 정답인지 아닌지 확인#k해야 합니다. 그럼 힘내주세요!")
			return
		end
		local answer = sm:get_property("stage8rand")
		local pos, total = area_pattern(map)
		if total ~= 5 then
			me:dialog(npc, "아직 5개의 정답 발판을 찾지 못하신것 같군요. 발판 끝에 아슬아슬하게 서계시지 말고 발판 가운데에 정화히 서 계셔야 정답 여부가 확인 가능합니다. 이점 주의해 주시기 바랍니다.")
			return
		end
		if pos == answer then
			map:clear_effect()
			sm:set_property("stage", "9")
			sm:set_property("guideRead", "0")
			give_exp(sm, 7200)
			me:dialog(npc, "다음 스테이지로 통하는 포탈이 열렸습니다.")
		else
			map:show_effect("quest/party/wrong_kor")
			map:play_sound("Party1/Failed")
		end
	end
}
