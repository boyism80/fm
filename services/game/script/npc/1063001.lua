-- NPC name (String.wz/Npc.img.xml): 파란색 꽃 무더기

function on_click(me, npc)
	me:map(105040300)
	local q = me:quest(2053)
	if q ~= nil and q:started() then
		local count = 0
		for _, it in pairs(me:item(4031026)) do
			count = count + it:count()
		end
		if count < 20 then
			local amt = math.random() < 0.3 and 20 or 15
			me:exchange(nil, { item = { [4031026] = amt } })
			return
		end
	end
	local ores = { 4020005, 4020006, 4020004, 4020001, 4020003, 4020000, 4020002 }
	local pick = ores[math.random(1, #ores)]
	me:exchange(nil, { item = { [pick] = 2 } })
end
