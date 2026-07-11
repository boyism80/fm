-- NPC name (String.wz/Npc.img.xml): 약초 덤불

function on_click(me, npc)
	local q = me:quest(2051)
	if q ~= nil and q:started() then
		local count = 0
		for _, it in pairs(me:item(4031032)) do
			count = count + it:count()
		end
		if count < 1 then
			me:exchange(nil, { item = { [4031032] = 1 } })
			me:map(101000000, 0)
			return
		end
	end
	local rand = math.random(1, 4)
	if rand == 1 then
		me:exchange(nil, { item = { [4020007] = 2 } })
	elseif rand == 2 then
		me:exchange(nil, { item = { [4020008] = 2 } })
	elseif rand == 3 then
		me:exchange(nil, { item = { [4010006] = 2 } })
	else
		me:exchange(nil, { item = { [1032013] = 1 } })
	end
	me:map(101000000, 0)
end
