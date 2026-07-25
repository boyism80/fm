-- NPC name (String.wz/Npc.img.xml): 아레다

return {
	on_click = function(me, npc)
		if not me:dialog_accept(npc, "거기 누구냐!") then
			return
		end
		me:notice("왕비에게 들켜 궁전에서 쫓겨났습니다.", Msg.PinkText)
		me:map(260000300)
	end
}
