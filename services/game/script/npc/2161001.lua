-- NPC name (String.wz/Npc.img.xml): 이피아

return {
	on_click = function(me, npc)
		if me:quest(3173):started() or me:quest(3175):started() then
			me:map(211070200)
			return
		end
		if me:quest(3178):started() then
			me:map(211070300)
			return
		end
		me:dialog(npc, "길이 막혀있습니다.")
	end
}
