-- NPC name (String.wz/Npc.img.xml): 가마

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "정말 이곳에서 나가 #b헤네시스#k 맵으로 돌아가고 싶나?") then
			return
		end
		me:map(100000000)
	end
}
