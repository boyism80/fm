-- NPC name (String.wz/Npc.img.xml): 바이올렛 벌룬

local pq = require("script/lib/party_quest")

local REPEAT_QUEST = 199600
local KEY_ID = 4001023
local BONUS_MAP = 922011000
local RANKING_QUEST = 1202

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
			pq.party_warp(sm, BONUS_MAP)
			return
		end
		local guide = sm:get_property("guideRead")
		if guide ~= "s" or not pq.is_leader(me) then
			if pq.is_leader(me) then
				sm:set_property("guideRead", "s")
			end
			me:dialog(npc, "안녕하세요. 드디어 여기까지 오셨군요... 이제 이 모든 소동을 일으킨 장본인을 쓰러뜨릴 시간입니다. 오른쪽으로 가보면 몬스터가 한 마리 있는데 쓰러뜨리면 엄청난 몸집의 #b알리샤르#k가 나타날 겁니다. 녀석은 여러분들 때문에 지금 몹시 화가 나 있는 상태이니 조심하세요.\r\n파티원이 모두 함께 녀석을 쓰러뜨리고 녀석이 지니고 있던 #b차원의 열쇠#k를 저에게 가져와 주시면 됩니다. 그 열쇠만 녀석에게서 빼앗는다면 다시 차원의 문이 열리는 일은 없겠지요. 그럼 여러분들만 믿겠습니다. 힘내 주세요!")
			return
		end
		if not pq.has_item(me, KEY_ID, 1) then
			me:dialog(npc, "아직 #b차원의 열쇠#k를 획득하지 못하신 모양이군요. 차원의 열쇠는 #b알리샤르#k를 쓰러뜨리면 얻으실 수 있습니다. 그럼 힘내 주세요!")
			return
		end
		map:show_effect("quest/party/clear")
		map:play_sound("Party1/Clear")
		pq.remove_all(KEY_ID, me)
		sm:set_property("stage", tostring(stage + 1))
		sm:set_property("guideRead", "0")
		give_exp(sm, 38500)
		for _, p in ipairs(sm:players()) do
			if p ~= nil then
				p:end_party_quest(RANKING_QUEST)
			end
		end
		me:dialog(npc, "차원의 열쇠를 가져 오셨군요! 여러분들은 모든 스테이지를 훌륭히 클리어 하셨습니다. 저에게 다시 말을 걸어주시면 파티원 전원이 보너스 맵으로 이동됩니다.")
	end
}
