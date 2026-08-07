-- NPC name (String.wz/Npc.img.xml): 레나리우

local alliance_create_ok = 0

local alliance_disband_ok = 0
local alliance_disband_not_in_alliance = 1
local alliance_disband_not_leader = 2
local alliance_disband_not_guild_master = 3
local alliance_disband_failed = 4

local alliance_inc_capacity_ok = 0
local alliance_inc_capacity_not_in_alliance = 1
local alliance_inc_capacity_not_leader = 2
local alliance_inc_capacity_not_guild_master = 3
local alliance_inc_capacity_capacity_max = 5
local alliance_inc_capacity_meso_cost = 5000000
local alliance_inc_capacity_failed = 6

local function try_create_alliance(npc, me)
	local party = me:party()
	if not party then
		me:dialog(npc, '길드장 두명이 파티를 맺고 있어야 길드 연합을 만들 수 있어요.')
		return
	end
	if party:leader_id() ~= me:id() then
		me:dialog(npc, '길드장 두명이 파티를 맺고 있어야 길드 연합을 만들 수 있어요.')
		return
	end
	local party_members = party:members()
	if #party_members ~= 2 then
		me:dialog(npc, '길드장 두명이 파티를 맺고 있어야 길드 연합을 만들 수 있어요.')
		return
	end
	local partner_id = nil
	for _, pm in ipairs(party_members) do
		if pm:id() ~= me:id() then
			partner_id = pm:id()
			break
		end
	end
	if not partner_id then
		me:dialog(npc, '길드장 두명이 파티를 맺고 있어야 길드 연합을 만들 수 있어요.')
		return
	end
	local map = me:map()
	if not map then
		me:dialog(npc, '파티원 두명이 이곳에 모여 잇는지 확인해 주세요.')
		return
	end
	local partner = map:characters()[partner_id]
	if not partner then
		me:dialog(npc, '파티원 두명이 이곳에 모여 잇는지 확인해 주세요.')
		return
	end
	local g = me:guild()
	if not g or g:rank(me) ~= 1 then
		me:dialog(npc, '길드장만 길드 연합을 등록할 수 있어요.')
		return
	end
	local partner_g = partner:guild()
	if not partner_g or partner_g:leader_id() ~= partner_id then
		me:dialog(npc, '다른 파티원이 길드장이 아닌 것 같아요.')
		return
	end
	if g:alliance() then
		me:dialog(npc, '이미 길드 연합에 가입되어 있으면 길드 연합을 만들 수 없습니다.')
		return
	end
	if partner_g:alliance() then
		me:dialog(npc, '다른 파티원이 이미 다른 길드 연합에 가입되어 있어요.')
		return
	end
	if not me:dialog_yes_no(npc, '정말 길드 연합을 만들고 싶으세요?') then
		return
	end
	local alliance_name = me:dialog_input(npc, '생성할 길드 연합 이름을 입력해 주세요. (최대 12바이트)')
	if alliance_name == nil or alliance_name == '' then
		return
	end
	if not me:dialog_yes_no(npc, string.format('입력하신 길드 연합 이름은 #b%s#k 입니다. 정말 이 이름으로 길드 연합을 만들고 싶으세요?', alliance_name)) then
		return
	end
	if me:meso() < 5000000 then
		me:dialog(npc, '사용할 수 없는 이름이거나 이미 존재하는 이름입니다. 다른 이름을 입력해 주세요. 혹은 수수료 500만메소가 부족하신건 아닌지도 확인해주세요.')
		return
	end
	local result = me:create_alliance(alliance_name)
	if result == alliance_create_ok then
		me:dialog(npc, '길드 연합이 새롭게 만들어졌습니다.')
	else
		me:dialog(npc, '사용할 수 없는 이름이거나 이미 존재하는 이름입니다. 다른 이름을 입력해 주세요. 혹은 수수료 500만메소가 부족하신건 아닌지도 확인해주세요.')
	end
end

local function try_disband_alliance(npc, me)
	local g = me:guild()
	if not g or not g:alliance() then
		me:dialog(npc, '길드 연합이 존재하지 않습니다.')
		return
	end
	if g:rank(me) ~= 1 then
		me:dialog(npc, '길드 연합장만 길드 연합을 해체할 수 있어요.')
		return
	end
	if not me:dialog_yes_no(npc, '길드 연합을 정말 해체하고 싶으세요? 신중하게 결정해 주시기 바랍니다.') then
		return
	end
	local result = me:disband_alliance()
	if result == alliance_disband_ok then
		me:dialog(npc, '길드 연합이 해체되었습니다.')
	elseif result == alliance_disband_not_guild_master or result == alliance_disband_not_leader then
		me:dialog(npc, '길드 연합장만 길드 연합을 해체할 수 있어요.')
	elseif result == alliance_disband_not_in_alliance then
		me:dialog(npc, '길드 연합이 존재하지 않습니다.')
	else
		me:dialog(npc, '길드 연합을 해체하지 못했습니다. 잠시 후 다시 시도해 주세요.')
	end
end

local function try_inc_alliance_capacity(npc, me)
	local g = me:guild()
	if not g or not g:alliance() then
		me:dialog(npc, '길드 연합이 존재하지 않습니다.')
		return
	end
	if g:rank(me) ~= 1 then
		me:dialog(npc, '길드 연합장만 길드 수를 늘릴 수 있어요.')
		return
	end
	if not me:dialog_yes_no(npc, '길드 수를 늘리는데에는 수수료 5백만 메소가 소비됩니다. 정말 만들고 싶으세요?') then
		return
	end
	if me:meso() < alliance_inc_capacity_meso_cost then
		me:dialog(npc, '최대 길드 연합의 길드 수는 5개 까지 늘릴 수 있습니다. 또는 수수료가 부족하신 건 아닌지 확인해 주세요.')
		return
	end
	local result = me:inc_alliance_capacity()
	if result == alliance_inc_capacity_ok then
		me:meso(me:meso() - alliance_inc_capacity_meso_cost)
		me:dialog(npc, '길드 연합의 길드수를 늘렸어요.')
	elseif result == alliance_inc_capacity_capacity_max then
		me:dialog(npc, '최대 길드 연합의 길드 수는 5개 까지 늘릴 수 있습니다. 또는 수수료가 부족하신 건 아닌지 확인해 주세요.')
	else
		me:dialog(npc, '길드 연합의 길드 수를 늘리지 못했습니다. 잠시 후 다시 시도해 주세요.')
	end
end

return {
	on_click = function(me, npc)
		local npc = 2010009
		local selected = me:dialog_list(npc,
			' 안녕하세요? #b레나리우#k라고 해요.',
			{
				'길드 연합이 무엇인지 알려주세요',
				'길드 연합을 만들려면 어떻게 해야 돼요?',
				'길드 연합을 만들고 싶어요.',
				'길드 연합의 길드 수를 늘리고 싶어요.',
				'길드 연합을 해체하고 싶어요.',
			})
		if selected == nil then
			return
		end
		if selected == 1 then
			me:dialog(npc, '여러개의 길드가 서로 모여서 만든 모임을 길드 연합이라고 해요. 저는 이렇게 만들어진 길드 연합을 관리하는 일을 하고 있답니다.')
		elseif selected == 2 then
			me:dialog(npc, '길드 연합을 만들려면 길드장 2명이 파티를 맺고 있어야 해요. 여기서 파티장이 길드 연합장이 된답니다.', false, true)
			me:dialog(npc, '2명의 길드장이 모였다면 500만 메소가 필요해요. 이건 길드 연합을 등록하는데 필요한 수수료에요.', false, true)
			me:dialog(npc, '그리고 또하나! 당연히 다른 길드 연합에 가입되어 있으면 새롭게 길드 연합을 만들지 못해요!', false, true)
		elseif selected == 3 then
			try_create_alliance(npc, me)
		elseif selected == 4 then
			try_inc_alliance_capacity(npc, me)
		elseif selected == 5 then
			try_disband_alliance(npc, me)
		end
	end
}
