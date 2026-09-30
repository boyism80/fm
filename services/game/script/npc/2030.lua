-- NPC name (String.wz/Npc.img.xml): 스테이지 도우미

local SCRIPT = "script/npc/2030.lua"
local MAX_MEMBERS = 4
local LOBBY = 123456771
local WAVES = {
	{ map = 123456772 },
	{ map = 123456773, mob = 3000016, x = -354, y = 67 },
	{ map = 123456774, mob = 3000017, x = -755, y = 388 },
	{ map = 123456775, mob = 3000019, x = -725, y = 476 },
	{ map = 123456776, mob = 3000014, x = 0, y = 0 },
}

return {
	spawn_wave = function(map, mob_id, x, y)
		map:spawn_mob(mob_id, x, y)
	end,

	on_click = function(me, npc)
		local party = me:party()
		if party == nil then
			me:dialog(npc, "파티를 구성하고 있지 않으시네요. 악용을 방지하여 이대로 퇴장됩니다.")
			me:map(LOBBY)
			return
		end
		if party:leader_id() ~= me:id() then
			me:dialog(npc, "파티장만이 저에게 말을 걸 수 있습니다.")
			return
		end
		local map = me:map()
		local characters = map:characters()
		local size = 0
		local in_map = 0
		for _, mem in pairs(party:members()) do
			size = size + 1
			if characters[mem:id()] ~= nil then
				if mem:class_id() == Class.GM then
					in_map = in_map + 2
				else
					in_map = in_map + 1
				end
			end
		end
		if size > MAX_MEMBERS or in_map < MAX_MEMBERS then
			me:dialog(npc, "다음 웨이브로 넘어가시기 위해선 파티원이 모두 파티장과 같은 맵에 있어야하고, 티어가 플래티넘이어야 합니다.")
			return
		end

		local map_id = map:wz():id()
		local index = nil
		for i, wave in ipairs(WAVES) do
			if wave.map == map_id then
				index = i
				break
			end
		end
		if index == nil then
			return
		end
		local dest = LOBBY
		if index < #WAVES then
			if next(map:mobs()) ~= nil then
				return
			end
			local wave = WAVES[index + 1]
			run_on_map(wave.map, SCRIPT, "spawn_wave", wave.mob, wave.x, wave.y)
			dest = wave.map
		end
		map:remove_npc(npc)
		local pid = party:id()
		for _, ch in pairs(characters) do
			local p = ch:party()
			if p ~= nil and p:id() == pid then
				ch:map(dest, 0)
			end
		end
	end
}
