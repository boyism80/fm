local dojo = require("script/lib/dojo")

return {
	on_enter = function(me)
		local sm = me:state_machine()
		local floor = dojo.current_floor(me)
		if sm == nil or floor == nil or not dojo.cleared(sm, floor) then
			me:message("아직 몬스터가 남아있습니다.", Msg.PinkText)
			return
		end
		me:show_effect(EffectType.QuestCompletion)
		local map = me:map()
		me:map(map, dojo.DOOR_PORTAL, { relocate = true })
		local door = map:reactor(dojo.DOOR_REACTOR)
		if door ~= nil then
			door:hit(1)
		end
	end
}
