-- NPC name (String.wz/Npc.img.xml): 덴마

local GROUP_NAME = "gojarani"

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm ~= nil and sm:group() ~= nil and sm:group():name() == GROUP_NAME then
			me:dialog(npc, "밖으로 나가려면 포탈을 이용하세요.")
			return
		end
		local group = state_machine(GROUP_NAME)
		if group == nil then
			me:dialog(npc, "지금은 입장할 수 없습니다.")
			return
		end
		local state = group:get_property("state")
		if state ~= nil and state ~= "" and state ~= "0" then
			me:dialog(npc, "이미 다른 사람이 안에 있습니다.")
			return
		end
		if not me:dialog_yes_no(npc, "입장하시겠습니까?") then
			return
		end
		local started, err = group:start_solo(me)
		if started == nil then
			me:dialog(npc, "지금은 입장할 수 없습니다.")
			if err ~= nil then
				log("gojarani start_solo:", err)
			end
		end
	end,
}
