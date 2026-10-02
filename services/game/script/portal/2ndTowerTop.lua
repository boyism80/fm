local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/2ndTowerTop.lua"
local TOWER = 211060401
local TEMP_KEY = 4032833

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
		local ok, count = run_on_map(TOWER, SCRIPT, "character_count")
		local allowed = me:quest(3140):started()
			and (pq.has_item(me, TEMP_KEY) or me:quest(3166):completed())
		if ok == false or count ~= 0 or allowed == false then
			me:message("두번째 탑의 임시열쇠가 없거나 이미 누군가가 안에서 퀘스트를 진행중입니다.", Msg.PinkText)
			return
		end

		run_on_map(TOWER, SCRIPT, "reset_map")
		if pq.has_item(me, TEMP_KEY) then
			me:rmitem(TEMP_KEY, 1)
		end
		me:map(TOWER, 0)
	end
}
