local SCRIPT = "script/portal/gotoAni.lua"
local ARENA = 211061100

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	on_enter = function(me)
		local ok, count = run_on_map(ARENA, SCRIPT, "character_count")
		local questing = me:quest(3142):started() or me:quest(3147):started()
		if ok == false or count ~= 0 or questing == false then
			me:message("길이 막혀있거나 이미 누군가 안에 아니와의 결투를 진행중입니다.", Msg.PinkText)
			return
		end
		me:open_npc(2161006)
	end
}
