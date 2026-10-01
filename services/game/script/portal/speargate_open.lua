-- Portal (old/scripts/portal/speargate_open.js): 기사의 홀

local gq = require("script/lib/guild_quest")

local SPEAR_MAP = 990000440
local SPEARS = { "spear1", "spear2", "spear3", "spear4" }

local function spears_lit(sm)
	local spear_map = sm:group():map(SPEAR_MAP)
	if spear_map == nil then
		return false
	end
	for _, name in ipairs(SPEARS) do
		local spear = spear_map:find_reactor_name(name)
		if spear == nil or spear:state() < 1 then
			return false
		end
	end
	return true
end

return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		local gate = me:map():find_reactor_name("speargate")
		local opened = gate ~= nil and gate:state() == 4
		if opened == false and spears_lit(sm) == false then
			me:message("지금은 포탈이 닫혀있습니다.", Msg.PinkText)
			return
		end
		gq.gain_gp_once(me, sm, "gainGP00", 20)
		me:play_portal_sound()
		me:map(990000401)
	end
}
