-- NPC name (String.wz/Npc.img.xml): 박 경관

local GROUP_NAME = "tokyo_2095"
local MIN_LEVEL = 120
local MAX_LEVEL = 250

return {
	on_click = function(me, npc)
		local level = me:level()
		if level < MIN_LEVEL or level > MAX_LEVEL then
			me:dialog(npc, "레벨이 맞지 않아 입장할 수 없습니다.")
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "지금은 입장할 수 없습니다.")
			return
		end
		local state = group:get_property("state")
		if state ~= nil and state ~= "" and state ~= "0" then
			me:dialog(npc, "이미 다른 사람이 도전 중입니다.")
			return
		end
		if not me:dialog_yes_no(npc, "공원 - 무서운 기계 군사에 입장하시겠습니까?") then
			return
		end
		local sm, err = group:start_solo(me)
		if sm == nil then
			me:dialog(npc, "지금은 입장할 수 없습니다.")
			if err ~= nil then
				log("tokyo_2095 start_solo:", err)
			end
		end
	end,
}
