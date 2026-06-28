-- Reactor name (Reactor.wz/2200001.img.xml): 922000020,922000021,22020300

function on_reactor_2200001(reactor)
	local trigger = reactor:trigger()
	if trigger == nil then
		return
	end
	trigger:notice("어딘가로 이동됩니다.")
	local roll = math.random(1, 100)
	local map_id
	if roll <= 30 then
		map_id = 922000020
	elseif roll <= 80 then
		map_id = 922000021
	else
		map_id = 220020300
	end
	trigger:map(map_id, 0)
end
