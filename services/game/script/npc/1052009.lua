-- NPC name (String.wz/Npc.img.xml): 보물상자

function on_click(me, npc)
	local q = me:quest(2056)
	if q ~= nil and q:started() then
		local count = 0
		for _, it in pairs(me:item(4031040)) do
			count = count + it:count()
		end
		if count < 1 then
			me:exchange(nil, { item = { [4031040] = 1 } })
			me:map(103000000, 0)
			return
		end
	end
	local ores = { 4020005, 4020006, 4020004, 4020001, 4020003, 4020000, 4020002 }
	local pick = ores[math.random(1, #ores)]
	me:exchange(nil, { item = { [pick] = 2 } })
	me:map(103000000, 0)
end
