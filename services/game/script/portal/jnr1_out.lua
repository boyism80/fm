return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil then
			return
		end
		if sm:get_property("stage1_way_clear") ~= "1" then
			local map = me:map()
			if map ~= nil then
				local pq = require("script/lib/party_quest")
				if pq.mob_count(map) > 0 then
					me:notice("아직 몬스터가 남아 있습니다.", Msg.PinkText)
					return
				end
			end
		end
		me:play_portal_sound()
		me:map(926110100)
	end
}
