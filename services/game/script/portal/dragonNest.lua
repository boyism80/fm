local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/dragonNest.lua"
local NEST = 240040611
local EXIT = 240040610
local EGG = 4001094
local GUARDIAN = 2081008

return {
	character_count = function(map)
		local n = 0
		for _, _ in pairs(map:characters()) do
			n = n + 1
		end
		return n
	end,

	prepare_nest = function(map)
		map:reset()
		for _, n in pairs(map:npcs()) do
			if n:id() == GUARDIAN then
				map:remove_npc(n)
			end
		end
	end,

	on_enter = function(me)
		if pq.has_item(me, EGG) == false then
			me:message("이 안에 들어가기 위해서는 <나인 스피릿의 알> 이 필요하다.", Msg.PinkText)
			return
		end
		local ok, count = run_on_map(NEST, SCRIPT, "character_count")
		if ok == false or count ~= 0 then
			me:message("이미 이 안에 들어가 누군가가 퀘스트를 완수하는 중입니다. 잠시 후 다시 시도해 주세요.", Msg.PinkText)
			return
		end
		run_on_map(NEST, SCRIPT, "prepare_nest")
		me:play_portal_sound()
		local function on_arrive(me)
			me:clock(300, function(me)
				me:map(EXIT)
			end)
		end
		me:map(NEST, 0, { callback = on_arrive })
	end
}
