local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/3rdTowerTop.lua"
local TOWER = 211060601
local TEMP_KEY = 4032834
local GUARD = 8840004
local GUARD_X = { 1483, 1683, 1880, 2347, 2108 }
local GUARD_Y = -139

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	prepare_tower = function(map)
		map:reset()
		for _, x in ipairs(GUARD_X) do
			map:spawn_mob(GUARD, x, GUARD_Y)
		end
	end,

	on_enter = function(me)
		local ok, count = run_on_map(TOWER, SCRIPT, "character_count")
		local allowed = me:quest(3141):started()
			and (pq.has_item(me, TEMP_KEY) or me:quest(3167):completed())
		if ok == false or count ~= 0 or allowed == false then
			me:message("세번째 탑의 임시열쇠가 없거나 이미 누군가가 안에서 퀘스트를 진행중입니다.", Msg.PinkText)
			return
		end

		run_on_map(TOWER, SCRIPT, "prepare_tower")
		if pq.has_item(me, TEMP_KEY) then
			me:rmitem(TEMP_KEY, 1)
		end
		me:map(TOWER, 0)
	end
}
