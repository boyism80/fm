-- NPC name (String.wz/Npc.img.xml): 클로이

return {
	on_click = function(me, npc)
		me:dialog(npc, "펫에 대해 궁금한 점이 있는가..?")
	end
}
