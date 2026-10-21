-- Map script (onUserEnter TD_EC_title): 입장 제목 효과

return {
	on_enter = function(me, map)
		me:map_effect("temaD/enter/event")
	end,
}
