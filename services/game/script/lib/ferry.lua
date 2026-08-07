local pq = require("script/lib/party_quest")

local M = {}

function M.prop(group_name, key)
	local g = state_machine(group_name)
	if g == nil then
		return ""
	end
	return g:get_property(key)
end

function M.ticket_for(me, ticket_low, ticket_high)
	if me:level() < 30 then
		return ticket_low
	end
	return ticket_high
end

function M.board(me, npc, opts)
	local group_name = opts.group
	local ticket = M.ticket_for(me, opts.ticket_low, opts.ticket_high)
	local warp_delta = opts.warp_delta or 1
	local min_level = opts.min_level or 0
	if min_level > 0 and me:level() < min_level then
		me:dialog(npc, opts.min_level_text or "레벨이 낮아 탑승할 수 없습니다.")
		return
	end
	local selected = me:dialog_list(npc,
		opts.intro or "배는 매 시간 정각 기준으로 10분 마다 출발하고 있으며, 출발 5분 전부터 표를 받고 있답니다.\r\n",
		{ opts.board_label or "배에 탑승하고 싶습니다." })
	if selected == nil then
		me:dialog(npc, "아직 이곳에서 볼일이 남으신 모양이지요?")
		return
	end
	local g = state_machine(group_name)
	if g == nil then
		me:dialog(npc, "배를 사용할 수 없습니다.")
		return
	end
	if M.prop(group_name, "ready") ~= "true" then
		me:dialog(npc, opts.ready_text)
		return
	end
	if M.prop(group_name, "entry") == "true" then
		if not me:dialog_yes_no(npc, "아직 배에 탑승할 여유가 있다고 합니다. 정말 배에 탑승하고 싶으세요?") then
			me:dialog(npc, "아직 이곳에서 볼일이 남으신 모양이지요?")
			return
		end
		if not pq.has_item(me, ticket, 1) then
			me:dialog(npc, "흐음.. #b#t" .. ticket .. "##k은 분명 제대로 갖고 계신건가요? 티켓이 없으시다면 왼쪽에서 티켓을 구매하실 수 있답니다.")
			return
		end
		if me:exchange({ item = { [ticket] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "흐음.. #b#t" .. ticket .. "##k은 분명 제대로 갖고 계신건가요? 티켓이 없으시다면 왼쪽에서 티켓을 구매하실 수 있답니다.")
			return
		end
		local map = me:map()
		if map == nil or map:wz() == nil then
			return
		end
		me:map(map:wz():id() + warp_delta)
		return
	end
	if M.prop(group_name, "docked") == "true" then
		me:dialog(npc, "이미 배가 출항 준비중에 있습니다. 죄송하지만 지금은 배에 탑승하실 수 없답니다. 표는 출발하기 1분 이전에만 받고 있답니다. 다음 배를 기다려보세요.")
		return
	end
	me:dialog(npc, opts.departed_text)
end

function M.sell(me, npc, opts)
	local ticket = M.ticket_for(me, opts.ticket_low, opts.ticket_high)
	local cost
	if me:level() < 30 then
		cost = opts.cost_low
	else
		cost = opts.cost_high
	end
	if opts.min_level and me:level() < opts.min_level then
		me:dialog(npc, opts.min_level_text or "레벨이 낮아 표를 살 수 없습니다.")
		return
	end
	local intro = opts.intro
		or ("배는 매 시간 정각 기준으로 10분 마다 출발하고 있으며, 출발 5분 전부터 표를 받고 있답니다.\r\n" ..
			"당신은 #b#t" .. ticket .. "##k이 필요하실 것 같군요. 요금은 #b" .. cost .. "#k 메소 입니다. 어떠세요? 구매해보시겠어요?")
	if not me:dialog_yes_no(npc, intro) then
		me:dialog(npc, opts.cancel_text or "그런가요? 새로운 대륙으로 모험을 떠나는것도 때로는 괜찮지 않을까요? 마음이 바뀌시면 다시 찾아오세요.")
		return
	end
	if me:exchange({ meso = cost }, { item = { [ticket] = 1 } }) ~= ExchangeResult.OK then
		me:dialog(npc, "흐음, 메소가 부족한건 아닌지, 인벤토리 공간이 부족한건 아닌지 다시 한번 확인해 주세요.")
		return
	end
	me:dialog(npc, "#b#t" .. ticket .. "##k은 잘 받으셨나요? 배 출발 시간에 늦지 않게 탑승해 주세요~")
end

function M.exit_waiting(me, npc, delta)
	if not me:dialog_yes_no(npc, "아직 배가 출발하려면 잠시 있어야 합니다만.. 지금 내려서 정거장으로 돌아가실 수도 있답니다. 어떠세요? 지금 돌아가보시겠어요? 지금 내려도 표는 환불되지 않으니 신중하게 생각해주세요.") then
		me:dialog(npc, "잠시 후 배가 출발할 예정이니 조금만 기다려 주세요~")
		return
	end
	local map = me:map()
	if map == nil or map:wz() == nil then
		return
	end
	me:map(map:wz():id() + (delta or -1))
end

return M
