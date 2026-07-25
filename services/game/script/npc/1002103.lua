-- NPC name (String.wz/Npc.img.xml): 리더 알

return {
	on_click = function(me, npc)
		me:dialog(npc, "패밀리에 대해 궁금한게 있나?", false, false)
	end
}
