-- Reactor name (Reactor.wz/2008006.img.xml): 전축

function on_reactor_2008006(reactor, item)
	local wz = item:wz()
	if wz == nil then
		return false
	end
	local item_id = wz:id()
	if item_id < 4001056 or item_id > 4001062 then
		return false
	end
	return item:count() == reactor:react_item_quantity()
end
