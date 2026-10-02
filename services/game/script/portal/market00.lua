local DEFAULT_RETURN = 102000000

local EXIT_PORTALS = {
	[100000100] = "market00",
	[102000000] = "market00",
	[120000200] = "market00",
	[200000000] = "market00",
	[211000100] = "market00",
	[220000000] = "market00",
	[221000000] = "market00",
	[222000000] = "market00",
	[230000000] = "market01",
	[240000000] = "market00",
	[250000000] = "market00",
	[251000000] = "market00",
	[260000000] = "market00",
	[500000000] = "market00",
}

return {
	on_enter = function(me)
		local return_map = me:saved_location("FREE_MARKET") or DEFAULT_RETURN
		me:clear_saved_location("FREE_MARKET")
		me:map(return_map, EXIT_PORTALS[return_map] or 0)
	end
}
