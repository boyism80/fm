function on_reactor_9201002(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:music('Bgm10/Eregos')
	local x, y = reactor:position()
	map:spawn_mob(9300028, x, y - 10)
	map:spawn_mob(9300031, 130, 90)
	map:spawn_mob(9300032, 540, 90)
	map:spawn_mob(9300029, 130, 150)
	map:spawn_mob(9300030, 540, 150)
end
