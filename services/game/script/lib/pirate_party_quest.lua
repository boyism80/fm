local M = {}

function M.lock_door(reactor, mob_id)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:set_respawn(false, mob_id)
	map:message("문이 잠겼습니다. 해적들이 더 이상 나오지 않습니다.")
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	local sm = trigger:state_machine()
	if sm == nil then
		return
	end
	local n = tonumber(sm:get_property("stage4")) or 0
	sm:set_property("stage4", tostring(n + 1))
end

return M
