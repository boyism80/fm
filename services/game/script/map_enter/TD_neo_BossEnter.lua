-- Map script (onUserEnter TD_neo_BossEnter): 네오시티 보스 방은 들어올 때마다 처음 상태로 되돌린다

return {
	on_enter = function(me, map)
		map:reset()
	end,
}
