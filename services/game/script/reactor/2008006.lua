-- Reactor name (Reactor.wz/2008006.img.xml): 전축

local MUSIC = {
	[4001056] = { path = "Bgm08/ForTheGlory", day = "1" },
	[4001057] = { path = "Bgm11/Aquarium", day = "2" },
	[4001058] = { path = "Bgm06/WelcomeToTheHell", day = "3" },
	[4001059] = { path = "Bgm06/FantasticThinking", day = "4" },
	[4001060] = { path = "Bgm02/EvilEyes", day = "5" },
	[4001061] = { path = "Bgm10/TheWayGrotesque", day = "6" },
	[4001062] = { path = "Bgm01/MoonlightShadow", day = "7" },
}

local function find_sm(map)
	if map == nil then
		return nil
	end
	for _, ch in pairs(map:characters()) do
		if ch ~= nil then
			local sm = ch:state_machine()
			if sm ~= nil then
				return sm
			end
		end
	end
	return nil
end

return {
	on_reactor = function(reactor, item)
		if item ~= nil then
			local wz = item:wz()
			if wz == nil then
				return false
			end
			local item_id = wz:id()
			if MUSIC[item_id] == nil then
				return false
			end
			if item:count() ~= reactor:react_item_quantity() then
				return false
			end
			local sm = find_sm(reactor:map())
			if sm ~= nil then
				sm:set_property("stage3_item", tostring(item_id))
			end
			return true
		end
		local map = reactor:map()
		local sm = find_sm(map)
		if sm == nil or map == nil then
			return
		end
		local item_id = tonumber(sm:get_property("stage3_item")) or 0
		local info = MUSIC[item_id]
		if info == nil then
			return
		end
		map:music(info.path)
		sm:set_property("stage3_music", info.day)
		reactor:hit(1)
	end
}
