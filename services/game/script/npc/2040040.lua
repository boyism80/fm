-- NPC name (String.wz/Npc.img.xml): 그린 벌룬

local pq = require("script/lib/party_quest")

local PASS_ID = 4001022
local REPEAT_QUEST = 199600
local COUNT = 24
local NEXT_STAGE = 7
local EXP = 5400
local GUIDE = "안녕하세요. 다섯번째 스테이지에 오신 것을 환영합니다. 이곳에는 여러 공간이 있고 그 안에서 몬스터를 쓰러뜨리거나 해서 역시 #b차원의 통행증 24장#k을 모아오시면 됩니다. 다만 특이한 점은 특정 직업이 아니고서는 #b차원의 통행증#k을 얻지 못하는 경우도 있으니 주의해야 한다는 것입니다.\r\n한 가지 힌트를 드리면 이 안에는 절대로 죽일 수 없는 #b차원의 킹 블록골렘#k이 있는데 도적이 아니고서는 녀석의 반대편으로 통과하기 힘들 것입니다. 또한 마법사가 아니고서는 갈 수 없는 곳도 있죠. 방법은 여러분들이 생각해 보세요. 그럼 힘내 주세요!"
local LACK = "아직 #b차원의 통행증 24장#k을 모으지 못하신 모양이지요? 이곳에는 여러 공간이 있고 그 안에서 몬스터를 쓰러뜨리거나 해서 역시 #b차원의 통행증 24장#k을 모아오시면 됩니다. 다만 특이한 점은 특정 직업이 아니고서는 #b차원의 통행증#k을 얻지 못하는 경우도 있으니 주의해야 한다는 것입니다.\r\n한 가지 힌트를 드리면 이 안에는 절대로 죽일 수 없는 #b차원의 킹 블록골렘#k이 있는데 도적이 아니고서는 녀석의 반대편으로 통과하기 힘들 것입니다. 또한 마법사가 아니고서는 갈 수 없는 곳도 있죠. 방법은 여러분들이 생각해 보세요."

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
