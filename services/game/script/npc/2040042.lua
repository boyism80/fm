-- NPC name (String.wz/Npc.img.xml): 스카이블루 벌룬

local pq = require("script/lib/party_quest")

local PASS_ID = 4001022
local REPEAT_QUEST = 199600
local COUNT = 3
local NEXT_STAGE = 8
local EXP = 6600
local GUIDE = "안녕하세요. 일곱번째 스테이지에 오신 것을 환영합니다. 이 곳에는 아주 아주 강력한 몬스터가 있습니다. 바로 #b차원의 롬바드#k라는 녀석이죠. 이 녀석을 쓰러뜨리면 다음 스테이지로 가는데 필요한 #b차원의 통행증#k을 줍니다. #b차원의 통행증 3장#k을 모아 파티장에게 주세요.\r\n녀석을 불러내는 방법은 멀리 떨어진 곳의 몬스터를 쓰러뜨리는 것입니다. 너무 멀어서 원거리 공격이 아니면 힘들겠지만... 아참... 차원의 롬바드는 보통 녀석이 아니니까 조심해 주세요. 얕보았다가는 큰 코 다칠수도 있으니까요. 그럼 힘내주세요!"
local LACK = "아직 #b차원의 통행증 3장#k을 모으지 못하신 모양이지요? #b차원의 롬바드#k를 쓰러뜨리고 #b차원의 통행증 3장#k을 모아와 주세요!"

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
		local cur = map_stage(wz:id())
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
