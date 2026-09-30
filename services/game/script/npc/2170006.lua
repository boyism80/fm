-- NPC name (String.wz/Npc.img.xml): 콜로세움

local SCRIPT = "script/npc/2170006.lua"
local BASE_MAP = 200101500
local ROOMS = 10

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	prepare_room = function(map)
		map:reset()
		map:spawn_mob(6160003, 642, 227)
	end,

	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "위험지역! 크세르크세스가 있는 곳입니다. 입장하시겠습니까?\r\n\r\n입장조건 : 레벨 50이상#b", {
			"입장 한다.",
			"되돌아 간다.",
		})
		if sel ~= 1 then
			return
		end
		local party = me:party()
		if party == nil then
			me:dialog(npc, "파티를 맺어야만 입장하실 수 있습니다.")
			return
		end
		if party:leader_id() ~= me:id() then
			me:dialog(npc, "파티장만 입장하실 수 있습니다.")
			return
		end
		for i = 0, ROOMS - 1 do
			local map_id = BASE_MAP + i
			if id2map(map_id) ~= nil then
				local ok, count = run_on_map(map_id, SCRIPT, "character_count")
				if ok and count == 0 then
					run_on_map(map_id, SCRIPT, "prepare_room")
					local pid = party:id()
					for _, ch in pairs(me:map():characters()) do
						local p = ch:party()
						if p ~= nil and p:id() == pid then
							ch:map(map_id, 1)
						end
					end
					return
				end
			end
		end
		me:dialog(npc, "모든 지역에 다른 파티가 입장 중입니다.")
	end
}
