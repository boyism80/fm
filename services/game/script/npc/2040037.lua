-- NPC name (String.wz/Npc.img.xml): 오렌지 벌룬

local pq = require("script/lib/party_quest")

local PASS_ID = 4001022
local REPEAT_QUEST = 199600
local COUNT = 15
local NEXT_STAGE = 3
local EXP = 3600
local GUIDE = "안녕하세요. 두번째 스테이지에 오신 것을 환영합니다. 맵 이곳 저곳을 돌아다니다 보면 몇 개의 상자를 찾으실 수 있을 겁니다. 상자를 부수면 다른 맵으로 이동되거나 차원의 통행증을 줄겁니다. 상자를 모두 조사해서 #b차원의 통행증 15장#k을 모아서 저에게 가져와 주시기 바랍니다. 차원의 통행증을 얻으면 파티장이 얻은 것을 모두 모아서 저에게 가져와 주시면 됩니다. 그럼 힘내주세요!"
local LACK = "아직 #b차원의 통행증 15장#k은 모으지 못하신 모양이지요? 맵 이곳 저곳에 있는 몬스터들을 잡고 나온 차원의 통행증을 파티장이 모아 제게 가져다 주시면 됩니다."

local function map_stage(map_id)
	if map_id == nil then
		return 0
	end
	return math.floor((map_id % 922010000) / 100)
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
		local wz = map:wz()
		if wz == nil then
			return
		end
		local stage_prop = sm:get_property("stage")
		if stage_prop == "" then
			sm:set_property("stage", "1")
			stage_prop = "1"
		end
		local stage = tonumber(stage_prop) or 1
		local cur = map_stage(wz.id)
		if stage > cur then
			me:dialog(npc, "포탈이 열렸어요~ 다음 스테이지로 이동해 주세요!")
			return
		end
		local guide = sm:get_property("guideRead")
		if guide ~= "s" or not pq.is_leader(me) then
			if pq.is_leader(me) then
				sm:set_property("guideRead", "s")
			end
			me:dialog(npc, GUIDE)
			return
		end
		if not pq.has_item(me, PASS_ID, COUNT) then
			me:dialog(npc, LACK)
			return
		end
		map:clear_effect()
		pq.remove_all(PASS_ID, me)
		sm:set_property("stage", tostring(NEXT_STAGE))
		sm:set_property("guideRead", "0")
		give_exp(sm, EXP)
		me:dialog(npc, "다음 스테이지로 통하는 포탈이 열렸습니다.")
	end
}
