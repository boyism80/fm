-- NPC name (String.wz/Npc.img.xml): 카르타

local pq = require("script/lib/party_quest")

local GROUP_NAME = "karta_cave"
local ENTRY_MAP = 230040001
local REQUIRED_ITEM = 4000175
local TARGET_ITEM = 4031472
local MIN_LEVEL = 120

return {
	on_click = function(me, npc)
		local quest = me:quest(6301)
		if not quest:started() then
			me:dialog(npc, "나는 #b바다 마녀 카르타#k다. 나를 귀찮게 했다가는 벌레로 만들어 버리는 수가 있으니 조심하도록 해.")
			return
		end
		if not pq.has_item(me, REQUIRED_ITEM) then
			me:dialog(npc, "#b#t4000175##k 없이는 일그러진 차원으로 들어갈 수 없다네.")
			return
		end
		if pq.has_item(me, TARGET_ITEM, 40) then
			me:dialog(npc, "이미 #b#t4031472##k을 40개 갖고 있는 것 같은데? 그렇다면 더 모을 필요 없지 않은가?")
			return
		end
		local party = me:party()
		if party == nil then
			me:dialog(npc, "파티가 없군. 파티를 만들고 내게 말을 걸게.")
			return
		end
		if party:leader_id() ~= me:id() then
			me:dialog(npc, "파티장이 입장 신청을 해야 하네.")
			return
		end
		for _, member in ipairs(party:members()) do
			local class_id = member:class_id()
			if class_id ~= 900 and class_id % 10 ~= 2 then
				me:dialog(npc, "흠.. 파티 중에 4차 전직을 하지 않은 플레이어가 있나?")
				return
			end
			if member:level() < MIN_LEVEL or member:channel_index() == nil or member:map_id() ~= ENTRY_MAP then
				me:dialog(npc, "파티원 모두 이곳에 모인 뒤 다시 말을 걸게.")
				return
			end
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "지금은 일그러진 차원으로 들어갈 수 없다네.")
			return
		end
		if group:get_property("started") == "true" then
			me:dialog(npc, "이미 다른 누군가가 다른 세계로 건너가서 임무를 수행 중인 것 같군. 나중에 다시 시도해 봐.")
			return
		end
		local sm, err = group:start_party(me, party)
		if sm == nil then
			me:dialog(npc, "지금은 일그러진 차원으로 들어갈 수 없다네.")
			if err ~= nil then
				log("karta_cave start_party:", err)
			end
		end
	end
}
