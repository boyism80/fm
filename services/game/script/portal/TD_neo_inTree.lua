local SCRIPT = "script/portal/TD_neo_inTree.lua"
local EXIT = 240070000
local QUESTS = {
	{ quest = 3719, map = 240070010 },
	{ quest = 3724, map = 240070020 },
	{ quest = 3730, map = 240070030 },
	{ quest = 3736, map = 240070040 },
	{ quest = 3742, map = 240070050 },
	{ quest = 3748, map = 240070060 },
}
local ROOMS = 10

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	reset_map = function(map)
		map:reset()
	end,

	on_enter = function(me)
		local base = nil
		for _, entry in ipairs(QUESTS) do
			if me:quest(entry.quest):started() then
				base = entry.map
			end
		end
		if base == nil then
			me:message("퀘스트를 진행 중일 때만 들어갈 수 있습니다.", Msg.PinkText)
			return
		end

		for i = 0, ROOMS - 1 do
			local room = base + i
			local ok, count = run_on_map(room, SCRIPT, "character_count")
			if ok and count == 0 then
				run_on_map(room, SCRIPT, "reset_map")
				me:play_portal_sound()
				local function on_arrive(me)
					me:clock(180, function(me)
						me:map(EXIT)
					end)
				end
				me:map(room, 0, { callback = on_arrive })
				return
			end
		end
		me:message("이미 다른 누군가가 들어가 있는 것 같다. 지금은 들어갈 수 없을 것 같다.", Msg.PinkText)
	end
}
