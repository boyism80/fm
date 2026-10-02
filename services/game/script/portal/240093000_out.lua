return {
	on_enter = function(me)
		if me:quest(31348):completed() and me:quest(31351):completed() == false then
			me:map(240092101)
		else
			me:map(240092100)
		end
	end
}
