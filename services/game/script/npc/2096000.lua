-- NPC name (String.wz/Npc.img.xml): 연습기록표

return {
	on_click = function(me, npc)
		for _, mob in pairs(me:map():mobs(5090001)) do
			mob:kill()
		end
	end
}
