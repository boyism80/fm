-- NPC name (String.wz/Npc.img.xml): 학

local ORBIS = 200000100
local ORBIS_STATION = 200000141
local MU_LUNG = 250000100
local HERB_TOWN = 251000000
local TO_MU_LUNG = 200090300
local TO_ORBIS = 200090310
local FLIGHT_SECONDS = 60
local FAR_FARE = 6000
local NEAR_FARE = 1500

local function fly(me, npc, fare, via, dest)
	if me:exchange({ meso = fare }, nil) ~= ExchangeResult.OK then
		me:dialog(npc, "충분한 메소를 소지하시지 않은 것 같군요.")
		return
	end
	me:map(via, 0)
	me:warp_later(dest, FLIGHT_SECONDS)
end

return {
	on_click = function(me, npc)
		local map_id = me:map():wz():id()
		if map_id == HERB_TOWN then
			if me:dialog_yes_no(npc, "안녕하세요? #b무릉도원#k으로 가고 싶으신가요? 지금 출발해 보시겠어요? 요금은 #b" .. NEAR_FARE .. " 메소#k입니다.") == false then
				me:dialog(npc, "마음이 바뀌면 다시 찾아오세요.")
				return
			end
			if me:exchange({ meso = NEAR_FARE }, nil) ~= ExchangeResult.OK then
				me:dialog(npc, "충분한 메소를 소지하시지 않은 것 같군요.")
				return
			end
			me:map(MU_LUNG)
		elseif map_id == ORBIS_STATION then
			local sel = me:dialog_list(npc, "안녕하세요? #b무릉도원#k으로 가고 싶으신가요? 지금 출발해 보시겠어요? 원하시는 것을 선택해 주세요.\r\n", {
				"#b무릉(" .. FAR_FARE .. " 메소)#k",
			})
			if sel == nil then
				return
			end
			fly(me, npc, FAR_FARE, TO_MU_LUNG, MU_LUNG)
		elseif map_id == MU_LUNG then
			local sel = me:dialog_list(npc, "안녕하세요? #b오르비스#k로 가고 싶으신가요? 지금 출발해 보시겠어요? 원하시는 것을 선택해 주세요.\r\n", {
				"#b오르비스(" .. FAR_FARE .. " 메소)#k",
				"#b백초마을(" .. NEAR_FARE .. " 메소)#k",
			})
			if sel == nil then
				return
			end
			if sel == 1 then
				fly(me, npc, FAR_FARE, TO_ORBIS, ORBIS)
				return
			end
			if me:dialog_yes_no(npc, "지금 #b백초마을#k 쪽으로 가보시겠어요? 요금은 #b" .. NEAR_FARE .. " 메소#k입니다.") == false then
				me:dialog(npc, "마음이 바뀌면 다시 찾아오세요.")
				return
			end
			if me:exchange({ meso = NEAR_FARE }, nil) ~= ExchangeResult.OK then
				me:dialog(npc, "충분한 메소르 소지하시지 않은 것 같군요.")
				return
			end
			me:map(HERB_TOWN)
		end
	end
}
