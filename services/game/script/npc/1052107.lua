-- NPC name (String.wz/Npc.img.xml): 작은 가로등

return {
	on_click = function(me, npc)
		for _, mob in pairs(me:map():mobs(5090000)) do
			mob:kill()
		end
	end
}
