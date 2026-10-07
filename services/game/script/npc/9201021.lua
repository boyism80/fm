-- NPC name (String.wz/Npc.img.xml): 비비안

local EXIT_MAP = 680000500

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:map(EXIT_MAP)
			return
		end
		local id = tostring(me:id())
		if id ~= sm:get_property("groom_id") and id ~= sm:get_property("bride_id") then
			me:dialog(npc, "한쌍의 커플이 드디어 맺어졌군요. 정말 아름다워요.")
			return
		end
		if sm:get_property("finale") == "0" then
			me:map(EXIT_MAP)
			return
		end
		if not me:dialog_yes_no(npc, "다음 이벤트로 넘어가고 싶으세요?") then
			me:dialog(npc, "아직 이곳에 더 머물고 싶으신건가요? 제한시간이 다되면 자동으로 넘어가니 그 전에 이동하고 싶다면 제게 다시 말을 걸어주세요.")
			return
		end
		sm:call_hook("on_finale")
	end
}
