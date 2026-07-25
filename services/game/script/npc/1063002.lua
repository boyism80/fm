-- NPC name (String.wz/Npc.img.xml): 흰색 꽃 무더기

return {
	on_click = function(me, npc)
		me:map(105040300)
		local q = me:quest(2054)
		if q ~= nil and q:started() then
			local count = 0
			for _, it in pairs(me:item(4031028)) do
				count = count + it:count()
			end
			if count < 30 then
				local amt = math.random() < 0.3 and 30 or 20
				me:exchange(nil, { item = { [4031028] = amt } })
				return
			end
		end
		local ores = { 4020007, 4020008, 4010006 }
		local pick = ores[math.random(1, #ores)]
		me:exchange(nil, { item = { [pick] = 2 } })
	end
}
