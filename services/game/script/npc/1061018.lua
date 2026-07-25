-- NPC name (String.wz/Npc.img.xml): 무영

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "안전한 곳으로 이동하시겠습니까?") then
			return
		end
		me:map(105100100)
	end
}
