return {
	on_enter = function(me)
		local sm = me:state_machine()
		if sm == nil or sm:get_property("stage5") ~= "0" then
			return
		end
		local map = me:map()
		if map == nil then
			return
		end
		local x, y = me:position()
		for id = 9300142, 9300146 do
			for _ = 1, 10 do
				map:spawn_mob(id, x, y)
			end
		end
		for _, n in pairs(map:npcs()) do
			if n ~= nil and n:id() == 2112000 then
				map:remove_npc(n)
			end
		end
		map:message("정체를 알 수 없는 과학자가 몬스터를 불러내고 황급히 사라졌다.")
		sm:set_property("stage5", "1")
	end
}
