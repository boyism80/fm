return {
	on_enter = function(me)
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		me:map(map:wz():id() + 2)
	end
}
