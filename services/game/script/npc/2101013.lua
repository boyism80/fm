-- NPC name (String.wz/Npc.img.xml): 카르카사

local TOWNS = { 100000000, 101000000, 102000000, 103000000, 104000000 }
local REFUSE = "음. 이 노선을 타시기에 조금 무서우신 모양이지요? 하지만 시간이 급하다면 얘기는 달라지겠지요."

return {
	on_click = function(me, npc)
		if not me:dialog_accept(npc, "이곳을 어떻게 찾으셨는지 모르겠군요. 하지만 제대로 찾아오셨다고 말씀드릴 수 있겠습니다. 이곳 니할 사막에서 빅토리아 아일랜드로 가는 직통 노선이 있거든요. 원하신다면 특별히 태워드리도록 하겠지만.. 정말 이 노선을 이용하고 싶으신가요?") then
			me:dialog(npc, REFUSE, false, true)
			return
		end
		if not me:dialog_accept(npc, "단, 명심하셔야 할 것이 두가지 있어요. 첫째로 이 노선은 화물 수송 노선으로 정식 노선이 아니기 때문에 #r어느 마을로 가게 될지는 장담할 수 없습니다.#k 둘째로 특별히 당신을 태워주는 거라서 돈이 좀 들어요. 수수료로 #b#e10000메소#k#n가 필요해요. 지금 당장 출발하는 수송선이 있는데 타시겠어요?") then
			me:dialog(npc, REFUSE, false, true)
			return
		end
		if not me:dialog(npc, "Okay, ready to takeoff~", false, true) then
			me:dialog(npc, REFUSE, false, true)
			return
		end

		local code = me:exchange({ meso = 10000 }, nil)
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "흐음. 메소가 부족하신건 아닌가요? 수수료는 #b#e10000#n#k메소니 돈을 충분히 갖고 있는지 확인해 보세요.", true, true)
			return
		end
		me:map(TOWNS[math.random(1, #TOWNS)])
	end
}
