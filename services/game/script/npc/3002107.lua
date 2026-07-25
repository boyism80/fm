-- NPC name (String.wz/Npc.img.xml): 랑

return {
	on_click = function(me, npc)
		me:dialog(npc, "안내 메뉴는 클릭할 수 없어요!", false, false)
	end
}
