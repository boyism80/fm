-- NPC name (String.wz/Npc.img.xml): 도의진

local MARBLE = 4001124
local SCROLLS = { 2043000, 2043100, 2043200, 2043300, 2043700, 2043800, 2044000, 2044100, 2044200, 2044300, 2044400, 2044500, 2044600, 2044700 }
local HERBS = { 4000276, 4000277, 4000278, 4000279, 4000280, 4000291, 4000292, 4000286, 4000287, 4000293, 4000294, 4000298, 4000284, 4000288, 4000285, 4000282, 4000295, 4000289, 4000296 }

return {
	on_click = function(me, npc)
		local npc_id = npc:id()

		local q = me:quest(3821)
		if q ~= nil and q:started() then
			q:force_complete(npc_id)
			me:dialog(npc, "퀘스트 완료.", false, true)
			return
		end

		local sel = me:dialog_list(npc, "다재다능한 나에게 무슨 볼일이신가? #b", {
			" 약을 만들고 싶어요",
			" 주문서를 만들고 싶어요",
			" 약재료를 기부하고 싶어요",
		})
		if sel == nil then
			return
		end

		if sel == 1 then
			local books = 0
			for _, it in pairs(me:item(4161030)) do
				books = books + it:count()
			end
			if books < 1 then
				me:dialog(npc, "약을 만들고 싶다면 먼저 약초와 약에 관련된 책을 먼저 공부하게나. 기본적인 지식이 있어야 무엇이든 할 수 있는 법이지.", false, true)
				return
			end
			me:dialog_list(npc, "어떤 약을 만들고 싶은가?\r\n#b", {
				" #i2022145:# #t2022145#",
				" #i2022146:# #t2022146#",
				" #i2022147:# #t2022147#",
				" #i2022148:# #t2022148#",
				" #i2022149:# #t2022149#",
				" #i2022150:# #t2022150#",
				" #i2050004:# #t2050004#",
				" #i4031554:# #t4031554#",
			})
			return
		end

		if sel == 2 then
			local options = {}
			for i, id in ipairs(SCROLLS) do
				options[i] = " #z" .. id .. "#"
			end
			local item_sel = me:dialog_list(npc, "만들고 싶은 주문서를 선택하게나.#b", options)
			if item_sel == nil then
				return
			end
			local item = SCROLLS[item_sel]
			local prompt = "흐음. #t" .. item .. "#를 만들어 보고 싶은가?  #t" .. item .. "#를 만들기 위해선 #b#t4001124##k과 #b강철의 원석 10개#k가 필요하다네."
			prompt = prompt .. "\r\n#i" .. MARBLE .. "# #t" .. MARBLE .. "# 100개"
			prompt = prompt .. "\r\n#i4010001# #t4010001# 10개"
			if not me:dialog_yes_no(npc, prompt) then
				return
			end
			local code = me:exchange({ item = { [MARBLE] = 100, [4010001] = 10 } }, { item = { [item] = 1 } })
			if code ~= ExchangeResult.OK then
				me:dialog(npc, "재료가 부족하거나 인벤토리에 공간이 부족한건 아닌지 확인해 보게나.")
			end
			return
		end

		local options = {}
		for i, id in ipairs(HERBS) do
			options[i] = " #z" .. id .. "#"
		end
		local herb_sel = me:dialog_list(npc, "엇? 그게 정말인가? 약재를 기부해주겠다니. 기부는 #b100#k개 단위로 받고 있다네. 기부를 해준다면 주문서를 만들 수 있는 구슬을 주겠네. 무엇을 기부해주겠는가?#b", options)
		if herb_sel == nil then
			return
		end
		local herb = HERBS[herb_sel]
		if not me:dialog_yes_no(npc, "정말 나에게 #b#t " .. herb .. "# 100개를 기부해 줄텐가?#k") then
			return
		end
		local code = me:exchange({ item = { [herb] = 100 } }, { item = { [MARBLE] = 1 } })
		if code == ExchangeResult.LackCapacity then
			me:dialog(npc, "흐음. 구슬을 받을 인벤토리 공간이 부족한건 아닌지 확인해 보게나.")
			return
		end
		if code ~= ExchangeResult.OK then
			me:dialog(npc, "재료가 부족하거나 인벤토리에 공간이 부족한건 아닌지 확인해 보게나.")
			return
		end
		me:dialog(npc, "고맙네~ 다음에도 부탁하겠네.")
	end
}
