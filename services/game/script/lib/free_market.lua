local M = {}

local FREE_MARKET = 910000000

function M.enter(me)
	if me:map():wz():id() == FREE_MARKET then
		return
	end
	if me:level() < 8 then
		me:message("레벨 8 이상만 자유시장에 입장할 수 있습니다.", Msg.PinkText)
		return
	end
	me:save_location("FREE_MARKET")
	me:play_portal_sound()
	me:map(FREE_MARKET, "st00")
end

return M
