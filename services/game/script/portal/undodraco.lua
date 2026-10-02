return {
	on_enter = function(me)
		if me:map():wz():id() == 200090500 then
			me:map(240000110)
		else
			me:map(270000100, 2)
		end
		me:morph(false)
	end
}
