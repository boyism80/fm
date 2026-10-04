local M = {}

function M.piece_hit(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local player = reactor:trigger()
	if player == nil then
		return
	end
	local sm = player:state_machine()
	if sm == nil then
		return
	end
	local status = tonumber(sm:get_property("status")) or 0
	sm:set_property("status", tostring(status + 1))
	local minerva = map:find_reactor_name("minerva")
	if minerva ~= nil then
		minerva:hit(minerva:state() + 1)
	end
end

return M
