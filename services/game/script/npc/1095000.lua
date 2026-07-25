-- NPC name (String.wz/Npc.img.xml): 슈린츠

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "아직 델리를 찾지 못한거야? 지금 그냥 나가고 싶어?") then
			return
		end
		me:map(120000104)
	end
}
