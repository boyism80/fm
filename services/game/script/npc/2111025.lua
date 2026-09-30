-- NPC name (String.wz/Npc.img.xml): 제어장치

return {
	on_click = function(me, npc)
		for _, mob in pairs(me:map():mobs(7090000)) do
			mob:kill()
		end
	end
}
