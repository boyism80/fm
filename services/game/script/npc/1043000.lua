-- NPC name (String.wz/Npc.img.xml): 꽃 무더기

return {
	on_click = function(me, npc)
		local q = me:quest(2050)
		if q ~= nil and q:started() then
			local count = 0
			for _, it in pairs(me:item(4031020)) do
				count = count + it:count()
			end
			if count < 1 then
				me:exchange(nil, { item = { [4031020] = 1 } })
				me:map(101000000, 0)
				return
			end
		end
		local ores = { 4020005, 4020006, 4020004, 4020001, 4020003, 4020000, 4020002 }
		local pick = ores[math.random(1, #ores)]
		me:exchange(nil, { item = { [pick] = 2 } })
		me:map(101000000, 0)
	end
}
