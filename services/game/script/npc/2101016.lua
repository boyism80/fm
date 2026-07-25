-- NPC name (String.wz/Npc.img.xml): 아레다

return {
	on_click = function(me, npc)
		me:dialog(npc, "테스트 엔피시 입니다!!", false, false)
	end
}
