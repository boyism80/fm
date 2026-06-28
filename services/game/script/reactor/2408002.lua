-- Reactor name (Reactor.wz/2408002.img.xml): 나무뿌리 구멍

function on_reactor_2408002(reactor, item)
	local wz = item:wz()
	if wz == nil then
		return false
	end
	local itemId = wz:id()
	if itemId < 4001088 or itemId > 4001091 then
		return false
	end
	return item:count() == reactor:react_item_quantity()
end
