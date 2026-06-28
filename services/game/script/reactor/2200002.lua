-- Reactor name (Reactor.wz/2200002.img.xml): 922010201로 추방

function on_reactor_2200002(reactor)
	local map = reactor:map()
	if map == nil then
		return
	end
	local trigger = reactor:trigger()
	local party = trigger ~= nil and trigger:party() or nil
	local party_id = party ~= nil and party:id() or nil
	for _, ch in pairs(map:characters()) do
		local ch_party = ch:party()
		if party_id ~= nil and ch_party ~= nil and ch_party:id() == party_id then
			ch:notice("함정에 빠져 어딘가로 이동됩니다.", Msg.PinkText)
			ch:map(922010201, 0)
		end
	end
end
