local M = {}

local ZAKUM_BODY_ID = 8800000
local ZAKUM_ARM_MIN = 8800003
local ZAKUM_ARM_MAX = 8800010

function M.on_arm_die(mob, attacker, map)
	if map == nil or mob == nil then
		return
	end
	local dying_oid = mob:oid()
	for _, m in pairs(map:mobs()) do
		if m:oid() ~= dying_oid then
			local id = m:id()
			if id >= ZAKUM_ARM_MIN and id <= ZAKUM_ARM_MAX then
				return
			end
		end
	end
	local spawn_x, spawn_y = mob:position()
	for _, m in pairs(map:mobs()) do
		if m:id() == ZAKUM_BODY_ID then
			spawn_x, spawn_y = m:position()
			break
		end
	end
	local oids = {}
	for _, m in pairs(map:mobs()) do
		oids[#oids + 1] = m:oid()
	end
	for _, oid in ipairs(oids) do
		map:remove_mob(oid, MobDieAnimation.FadeOut)
	end
	map:spawn_mob(ZAKUM_BODY_ID, spawn_x, spawn_y, -2)
end

return M
