function on_reactor_2401100(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	map:music('Bgm14/HonTale')
	map:spawn_mob(8810130, 71, 260)
	map:message('동굴이 울리면서 카오스 혼테일이 나타났습니다.')
end
