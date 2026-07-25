-- NPC name (String.wz/Npc.img.xml): 보물상자

return {
	on_click = function(me, npc)
		local q = me:quest(2057)
		if q ~= nil and q:started() then
			local count = 0
			for _, it in pairs(me:item(4031041)) do
				count = count + it:count()
			end
			if count < 1 then
				me:exchange(nil, { item = { [4031041] = 1 } })
				me:map(103000000, 0)
				return
			end
		end
		local rand = math.random(1, 3)
		if rand == 1 then
			me:exchange(nil, { item = { [4020007] = 2 } })
		elseif rand == 2 then
			me:exchange(nil, { item = { [4020008] = 2 } })
		else
			me:exchange(nil, { item = { [4010006] = 2 } })
		end
		me:map(103000000, 0)
	end
}
