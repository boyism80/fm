-- NPC name (String.wz/Npc.img.xml): 분홍색 꽃 무더기

return {
	on_click = function(me, npc)
		me:map(105040300)
		local q = me:quest(2052)
		if q ~= nil and q:started() then
			local count = 0
			for _, it in pairs(me:item(4031025)) do
				count = count + it:count()
			end
			if count < 10 then
				local amt = math.random() < 0.3 and 10 or 8
				me:exchange(nil, { item = { [4031025] = amt } })
				return
			end
		end
		local ores = { 4010003, 4010000, 4010002, 4010005, 4010004, 4010001 }
		local pick = ores[math.random(1, #ores)]
		me:exchange(nil, { item = { [pick] = 2 } })
	end
}
