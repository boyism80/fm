-- NPC name (String.wz/Npc.img.xml): 루이스

return {
	on_click = function(me, npc)
		if not me:dialog_yes_no(npc, "정말 이곳에서 나가고 #b엘리니아#k로 돌아가고 싶어?") then
			return
		end
		me:map(101000000, 0)
	end
}
