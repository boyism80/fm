-- NPC name (String.wz/Npc.img.xml): 비호 석상

local pq = require("script/lib/party_quest")
local gq = require("script/lib/guild_quest")

local RUBIAN = 4001024

local function clear_gp(sm)
	local min = (gq.WAIT_MS + gq.DURATION_MS - sm:time_left()) / 60000
	if min < 25 then
		return 850
	end
	if min < 30 then
		return math.floor(800 - (min - 25) * 20)
	end
	if min < 40 then
		return math.floor(580 - (min - 30) * 10)
	end
	if min < 60 then
		return math.floor(400 - (min - 40) * 5)
	end
	if min < 90 then
		return math.floor(260 - (min - 60) * 2)
	end
	return 200
end

return {
	on_click = function(me, npc)
		local sm = me:state_machine()
		if sm == nil then
			me:map(gq.EXIT_MAP)
			return
		end
		if gq.is_leader(me, sm) == false then
			return
		end
		if pq.has_item(me, RUBIAN) == false then
			me:dialog(npc, "루비안의 힘을 지니고 있는가?")
			return
		end

		pq.remove_all(RUBIAN, me)
		if sm:get_property("bossclear") == "" then
			sm:message("길드 대항전을 클리어 하였습니다. 누리스를 통해 이 맵을 나갈 수 있습니다.")
			gq.gain_gp_once(me, sm, "bossclear", clear_gp(sm))
		end
		sm:call_hook("on_clear")
	end
}
