-- NPC name (String.wz/Npc.img.xml): 삼장법사

local TARGET_MOB = 8800002
local MAX_HP = 999999999

return {
	on_click = function(me, npc)
		if me:role() ~= ROLE.Admin then
			return
		end

		local hp = tonumber(me:dialog_input(npc, "hp"))
		if hp == nil then
			return
		end
		hp = math.floor(math.max(0, math.min(hp, MAX_HP)))

		for _, mob in pairs(me:map():mobs(TARGET_MOB)) do
			mob:hp(hp)
		end
	end
}
