-- NPC name (String.wz/Npc.img.xml): 이프

local tickets_high = { 4031047, 4031074, 4031331, 4031576 }
local tickets_low = { 4031046, 4031073, 4031330, 4031575 }
local costs_high = { 3000, 2000, 2000, 2000 }
local costs_low = { 1500, 1000, 1000, 1000 }
local intervals = { 15, 10, 10, 10 }
local map_names = {
	"빅토리아 아일랜드의 엘리니아",
	"루디브리엄",
	"미나르숲의 리프레",
	"니할 사막의 아리안트",
}

return {
	on_click = function(me, npc)
		local labels = {}
		for i, name in ipairs(map_names) do
			labels[i] = name
		end
		local select = me:dialog_list(npc,
			"안녕하세요. 어느곳으로 가고 싶으세요? 원하는 곳에 가시려면 여기서 표를 구입하셔야 한답니다.",
			labels)
		if select == nil then
			me:dialog(npc, "그런가요? 여러곳으로 여행하는건 즐거운 일이지만 아직 볼일이 남아있으신가 보죠? 마음이 바뀌시면 다시 찾아오세요.")
			return
		end
		local idx = select + 1
		local low = me:level() < 30
		local ticket
		local cost
		if low then
			ticket = tickets_low[idx]
			cost = costs_low[idx]
		else
			ticket = tickets_high[idx]
			cost = costs_high[idx]
		end
		if not me:dialog_yes_no(npc,
			map_names[idx] .. " 로 가는 배는 매 " .. intervals[idx] .. " 분 마다 출발하고 있으며, 요금은 #b" .. cost .. " 메소#k입니다. 정말 #b#t" .. ticket .. "##k을 구입하시겠어요?") then
			me:dialog(npc, "그런가요? 여러곳으로 여행하는건 즐거운 일이지만 아직 볼일이 남아있으신가 보죠? 마음이 바뀌시면 다시 찾아오세요.")
			return
		end
		if me:exchange({ meso = cost }, { item = { [ticket] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "확실히 #b" .. cost .. " 메소#k를 잘 가지고 계신건가요? 아니면 인벤토리가 가득찬건 아닌지 확인해 주세요.")
			return
		end
		me:dialog(npc, "#b#t" .. ticket .. "##k은 잘 받으셨나요? 배는 오른쪽 포탈을 통해 정거장으로 가시면 탑승하실 수 있답니다. 배 출발 시간에 늦지 않게 탑승해주세요~")
	end,
}
