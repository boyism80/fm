-- Portal (old/scripts/portal/Tokyoboss2_go.js): 802000309 → 802000310

return {
	on_enter = function(me)
		me:play_portal_sound()
		me:map(802000310, 2)
	end,
}
