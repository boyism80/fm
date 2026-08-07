local rj = require("script/lib/romeo_juliet")

return {
	on_enter = function(me)
		rj.strip_items(me)
		me:play_portal_sound()
		local map = me:map()
		if map ~= nil and map:wz() ~= nil and map:wz():id() == 926100700 then
			me:map(261000011)
		else
			me:map(261000021)
		end
	end
}
