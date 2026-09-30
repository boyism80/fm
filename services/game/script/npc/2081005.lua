-- NPC name (String.wz/Npc.img.xml): 케로벤

local DRAGON = 2210003
local SCRIPT = "script/npc/2081005.lua"

local function spawn_of(map_id, name)
	if id2map(map_id) == nil then
		return 0
	end
	local ok, id = run_on_map(map_id, SCRIPT, "portal_id", name)
	if ok and id ~= nil then
		return id
	end
	return 0
end

return {
	portal_id = function(map, name)
		local portal = map:portal(name)
		if portal == nil then
			return nil
		end
		return portal:id()
	end,

	on_click = function(me, npc)
		local dragon = me:morph() == DRAGON
		local text = "인간이군! 내가 있는 한 이 곳에서 한걸음도 더 나아갈 수 없다. 썩 사라지거라!"
		if dragon then
			text = "오, 우리 동족이로군. 인간의 침입은 걱정 말라고. 내가 단단히 지키고 있으니까 말이야. 그럼 안으로 들어가게나."
		end
		if me:dialog(npc, text, false, true) == false then
			return
		end
		if dragon then
			me:map(240050000, spawn_of(240050000, "out00"))
			me:morph(false)
			return
		end
		local damage = 500
		if me:hp() < 500 then
			damage = me:hp() - 1
			if damage == me:hp() then
				damage = 0
			end
		end
		if damage > 0 then
			me:hp(me:hp() - damage)
		end
		me:map(240040600, spawn_of(240040600, "st00"))
	end
}
