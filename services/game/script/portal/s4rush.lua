local GROUP_NAME = "warrior_rush"
local QUEST = 6110
local PARTY_SIZE = 2

local function warrior(class_id)
	return class_id == 112 or class_id == 122 or class_id == 132 or class_id == 900
end

return {
	on_enter = function(me)
		if not me:quest(QUEST):started() then
			me:message("아직 이 시험에 들어갈 수 없습니다.")
			return
		end
		local party = me:party()
		if party == nil or party:leader_id() ~= me:id() then
			me:message("두 명의 전사로 파티를 구성하고 파티장이 입장해야 합니다.")
			return
		end
		local members = party:members()
		if #members ~= PARTY_SIZE then
			me:message("정확히 두 명의 전사가 필요합니다.")
			return
		end
		local characters = me:map():characters()
		for _, member in ipairs(members) do
			local character = characters[member:id()]
			if character == nil or not warrior(member:class_id()) or not character:quest(QUEST):started() then
				me:message("시험을 받을 두 전사가 모두 이곳에 있어야 합니다.")
				return
			end
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:message("지금은 시험장에 들어갈 수 없습니다.")
			return
		end
		local sm, err = group:start_party(me, party)
		if sm == nil then
			me:message("다른 파티가 시험 중입니다. 잠시 후 다시 시도하세요.")
			if err ~= nil then
				log("warrior_rush start_party:", err)
			end
		end
	end,
}
