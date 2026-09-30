-- NPC name (String.wz/Npc.img.xml): 앉아있는 여행객

return {
	on_click = function(me, npc)
		local map = me:map()
		local here = map ~= nil and map:template_id() == 450004750
		local ask = "그냥 같이 책이나 읽지..."
		if here then
			ask = "이 곳에서 나가겠어요? 같이 책이나 보시겠어요?"
		end
		if not me:dialog_yes_no(npc, ask) then
			return
		end
		if id2map(123456789) == nil then
			me:dialog(npc, "아직 갈 수 없는 곳입니다.")
			return
		end
		me:map(123456789)
	end
}
