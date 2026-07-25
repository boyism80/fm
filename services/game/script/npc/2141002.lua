-- NPC name (String.wz/Npc.img.xml): 잊혀진 신전관리인

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "정말 이곳에서 나가 안전한 곳으로 돌아가고 싶으십니까?") then
			return
		end
		me:map(270050300)
	end
}
