-- NPC name (String.wz/Npc.img.xml): 엄마 젖소

local COW = "1"

return {
	on_click = function(me, npc)
		local milk_q = me:quest(126640)
		local cow_q = me:quest(126641)
		if not milk_q:started() then
			milk_q:start("0")
		end
		if not cow_q:started() then
			cow_q:start("0")
		end

		local milk = milk_q:record()
		if milk == "0" then
			milk_q:record("1")
			cow_q:record(COW)
			me:dialog(npc, "우유를 적당히 채웁니다. 우유가 우유통의 1/3 정도 차있습니다.")
			return
		end
		if milk ~= "1" and milk ~= "2" then
			me:dialog(npc, "이미 우유통에 우유가 가득 차있습니다.")
			return
		end
		if cow_q:record() == COW then
			me:dialog(npc, "우유가 더 이상 나오지 않습니다.")
			return
		end

		if milk == "1" then
			milk_q:record("2")
			cow_q:record(COW)
			me:dialog(npc, "우유를 적당히 채웁니다. 우유가 우유통의 2/3 정도 차있습니다.")
			return
		end

		local code = me:exchange(nil, { item = { [4031850] = 1 } })
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "인벤토리 공간이 부족합니다.")
			return
		end
		milk_q:record("3")
		cow_q:record(COW)
		me:dialog(npc, "우유를 적당히 채웁니다. 우유통에 우유가 가득 차있습니다.")
	end
}
