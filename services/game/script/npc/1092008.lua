-- NPC name (String.wz/Npc.img.xml): 슈린츠

local QUEST = 6410
local PROGRESS_QUEST = 6411

return {
	on_click = function(me, npc)
		if me:quest(QUEST):started() == false then
			return
		end
		if me:dialog_yes_no(npc, "그럼 지금 바로 델리를 구하러 가볼래?") == false then
			return
		end
		if me:quest(PROGRESS_QUEST):record() ~= "" then
			me:dialog(npc, "이미 델리를 구하고 온 것 같은데? 다시 갈 필요는 없잖아?")
			return
		end
		local sm, err = state_machine("protect_delli"):start_solo(me)
		if sm == nil then
			log("protect_delli start_solo:", err)
			me:dialog(npc, "이미 이 안에 다른 누군가가 델리를 구출하는 중이야. 나중에 다시 시도해 봐.")
		end
	end
}
