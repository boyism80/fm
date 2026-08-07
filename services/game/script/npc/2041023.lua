-- NPC name (String.wz/Npc.img.xml): 플로

local pq = require("script/lib/party_quest")

local GROUP_NAME = "element_thanatos"
local ENTRY_MAP = 220050300
local PARTY_SIZE = 2
local MIN_LEVEL = 120

local function quest_ready(me)
	local fire_main = me:quest(6225)
	local fire_prerequisite = me:quest(6226)
	local ice_main = me:quest(6315)
	local ice_prerequisite = me:quest(6316)
	local fire_ready = fire_main ~= nil
		and fire_main:started()
		and fire_prerequisite ~= nil
		and fire_prerequisite:completed()
	local ice_ready = ice_main ~= nil
		and ice_main:started()
		and ice_prerequisite ~= nil
		and ice_prerequisite:completed()

	return fire_ready or ice_ready
end

return {
	on_click = function(me, npc)
		if not quest_ready(me) then
			me:dialog(npc, "속성의 타나토스요..? 지금 당신에겐 그것을 만날 필요 없어 보이는군요.")
			return
		end
		local party = me:party()
		if party == nil then
			me:dialog(npc, "먼저 파티를 구성하고 말을 걸어주세요.")
			return
		end
		local members = party:members()
		local count = 0
		for _, member in pairs(members) do
			if member ~= nil then
				count = count + 1
			end
		end
		local gm_solo = pq.is_gm(me) and count == 1
		if count ~= PARTY_SIZE and not gm_solo then
			me:dialog(npc, "파티 인원을 두명으로 맞춰주세요.")
			return
		end
		for _, member in pairs(members) do
			local class_id = member:class_id()
			if class_id ~= 212 and class_id ~= 222 and class_id ~= 900 then
				me:dialog(npc, "파티원 중 한명의 직업이 다른 세계로 입장하는데 맞지 않는 것 같군요.")
				return
			end
			if member:level() < MIN_LEVEL then
				me:dialog(npc, "파티원 중 한명의 레벨이 다른 세계로 입장하는데 맞지 않는 것 같군요.")
				return
			end
			if member:channel_index() == nil or member:map_id() ~= ENTRY_MAP then
				me:dialog(npc, "파티원 두 명 모두 이곳에 모인 뒤 다시 말을 걸어주세요.")
				return
			end
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "오류가 발생했습니다.")
			return
		end
		local state = group:get_property("state")
		if state ~= nil and state ~= "" and state ~= "0" then
			me:dialog(npc, "이미 다른 파티가 속성의 타나토스에 도전하고 있습니다.")
			return
		end
		local sm, err = group:start_party(me, party)
		if sm == nil then
			me:dialog(npc, "오류가 발생했습니다.")
			if err ~= nil then
				log("element_thanatos start_party:", err)
			end
		end
	end
}
