-- NPC name (String.wz/Npc.img.xml): 의미없는 존재

local GROUP_NAME = "archer_snipe"
local ENTRY_MAP = 105090200
local PARTY_SIZE = 2
local MIN_LEVEL = 120

return {
	on_click = function(me, npc)
		local quest = me:quest(6108)
		local prerequisite = me:quest(6107)
		if not quest:started() or not prerequisite:completed() then
			me:dialog(npc, ".....")
			return
		end
		local party = me:party()
		if party == nil then
			me:dialog(npc, "파티를 구성한 후 내게 말을 걸게.")
			return
		end
		if party:leader_id() ~= me:id() then
			me:dialog(npc, "4차 전직을 한 궁수 두 명이 파티를 구성해서 파티장이 내게 말을 걸어야 하네.")
			return
		end
		local members = party:members()
		if #members ~= PARTY_SIZE then
			me:dialog(npc, "파티 인원을 두 명으로 구성한 후 내게 말을 걸게.")
			return
		end
		for _, member in ipairs(members) do
			local class_id = member:class_id()
			if class_id ~= 312 and class_id ~= 322 and class_id ~= 900 then
				me:dialog(npc, "4차 전직을 한 궁수 두 명이 파티를 구성해서 파티장이 내게 말을 걸어야 하네.")
				return
			end
			if member:level() < MIN_LEVEL or member:channel_index() == nil or member:map_id() ~= ENTRY_MAP then
				me:dialog(npc, "4차 전직을 한 궁수 두 명이 이곳에 모인 뒤 다시 말을 걸어야 하네.")
				return
			end
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "지금은 다른 세계로 들어갈 수 없네.")
			return
		end
		if group:get_property("started") == "true" then
			me:dialog(npc, "이미 다른 누군가가 다른 세계로 건너가서 임무를 수행 중인 것 같군. 나중에 다시 시도해 보게.")
			return
		end
		local sm, err = group:start_party(me, party)
		if sm == nil then
			me:dialog(npc, "지금은 다른 세계로 들어갈 수 없네.")
			if err ~= nil then
				log("archer_snipe start_party:", err)
			end
		end
	end
}
