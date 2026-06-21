-- Reactor name (Reactor.wz/2401000.img.xml): 봉인된나인스피릿

function on_reactor_hit_2401000(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end

	map:music('Bgm14/HonTale')
	map:spawn_mob(8810026, 71, 260)
	map:message('동굴이 울리면서 혼테일이 나타났습니다.')
end
