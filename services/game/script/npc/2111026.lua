-- NPC name (String.wz/Npc.img.xml): 미완성 마법진

return {
	on_click = function(me, npc)
		for _, mob in pairs(me:map():mobs(8090000)) do
			mob:kill()
		end
	end
}
