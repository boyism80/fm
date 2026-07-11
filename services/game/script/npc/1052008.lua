-- NPC name (String.wz/Npc.img.xml): 보물상자

function on_click(me, npc)
	local q = me:quest(2055)
	if q ~= nil and q:started() then
		local count = 0
		for _, it in pairs(me:item(4031039)) do
			count = count + it:count()
		end
		if count < 1 then
			me:exchange(nil, { item = { [4031039] = 1 } })
			me:map(103000000, 0)
			return
		end
	end
	local ores = { 4010003, 4010000, 4010002, 4010005, 4010004, 4010001 }
	local pick = ores[math.random(1, #ores)]
	me:exchange(nil, { item = { [pick] = 2 } })
	me:map(103000000, 0)
end
