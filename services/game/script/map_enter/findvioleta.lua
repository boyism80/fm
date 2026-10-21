-- Map script (onUserEnter findvioleta): 바이올레타를 찾는 방은 들어올 때마다 처음 상태로 되돌린다

return {
	on_enter = function(me, map)
		map:reset()
	end,
}
