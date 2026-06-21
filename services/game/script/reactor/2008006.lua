-- Reactor name (Reactor.wz/2008006.img.xml): 전축

function on_item_drop_match_2008006(reactor, item)
	local wz = item:wz()
	if wz == nil then
		return false
	end
	local itemId = wz:id()
	if itemId < 4001056 or itemId > 4001062 then
		return false
	end
	return item:count() == reactor:react_item_quantity()
end
