local SCRIPT = "script/portal/240093200_in.lua"
local ROOM = 240093300

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	on_enter = function(me)
		local ok, count = run_on_map(ROOM, SCRIPT, "character_count")
		if me:quest(31351):started() == false or ok == false or count ~= 0 then
			me:message("길이 막혀있거나 이미 누군가 입장해 있습니다.", Msg.PinkText)
			return
		end
		me:open_npc(2210014)
	end
}
