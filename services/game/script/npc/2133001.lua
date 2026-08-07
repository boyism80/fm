-- NPC name (String.wz/Npc.img.xml): 엘린

local pq = require("script/lib/party_quest")

local PURPLE_STONE = 4001163
local MONSTER_MARBLE = 4001169
local PURIFY_BEAD = 2270004
local EXIT_MAP = 930000800

local function strip_pq_items(me)
	pq.remove_all(PURPLE_STONE, me)
	pq.remove_all(MONSTER_MARBLE, me)
	pq.remove_all(PURIFY_BEAD, me)
end

return {
	on_click = function(me, npc)
		local map = me:map()
		if map == nil then
			return
		end
		local wz = map:wz()
		if wz == nil then
			return
		end
		local map_id = wz.id
		if map_id == 930000000 then
			me:dialog(npc, "어서와, 중앙에 보이는 포탈에 들어가면 변신 마법을 걸어줄게.")
		elseif map_id == 930000010 then
			me:dialog(npc, "자신이 누군지 잘 확인하고 자신의 모습을 잊어버리지 않도록 조심해.")
		elseif map_id == 930000100 then
			me:dialog(npc, "독에 변질된 모든 몬스터들을 없애!")
		elseif map_id == 930000200 then
			me:dialog(npc, "중앙에서 희석된 독으로 가시 덤불을 없애고 진행해줘.")
		elseif map_id == 930000300 then
			local sm = me:state_machine()
			if sm == nil then
				return
			end
			map:clear_effect()
			pq.party_warp(sm, 930000400)
		elseif map_id == 930000400 then
			if pq.has_item(me, MONSTER_MARBLE, 20) and pq.is_leader(me) then
				local sm = me:state_machine()
				if sm == nil then
					return
				end
				map:clear_effect()
				pq.remove_all(MONSTER_MARBLE, me)
				pq.party_warp(sm, 930000500)
			elseif not pq.has_item(me, PURIFY_BEAD) and pq.is_leader(me) then
				local code = me:exchange({}, { item = { [PURIFY_BEAD] = 10 } })
				if code == ExchangeResult.LackCapacity then
					me:dialog(npc, "인벤토리 공간을 확보한 뒤 다시 말을 걸어줘.")
					return
				end
				if code == ExchangeResult.OK then
					me:dialog(npc, "여기 정화의 구슬을 줄게.")
				end
			else
				me:dialog(npc, "나에게 정화의 구슬을 받은 다음, 몬스터들을 캐치해서 몬스터 구슬 20개를 파티장이 가져와!")
			end
		elseif map_id == 930000600 then
			me:dialog(npc, "이제 괴인의 제단 위에 보라색 마력석을 올려놓아 봐.")
		elseif map_id == 930000700 then
			strip_pq_items(me)
			me:map(EXIT_MAP)
		end
	end
}
