-- NPC name (String.wz/Npc.img.xml): 헤라클

local guild_create_ok = 0
local guild_create_not_allowed = 1
local guild_create_insufficient_meso = 2
local guild_create_already_in_guild = 3
local guild_create_failed = 4

local guild_disband_ok = 0
local guild_disband_not_in_guild = 1
local guild_disband_not_master = 2
local guild_disband_failed = 3

local guild_capacity_increase_ok = 0
local guild_capacity_increase_not_in_guild = 1
local guild_capacity_increase_not_master = 2
local guild_capacity_increase_insufficient_meso = 3
local guild_capacity_increase_insufficient_gp = 4
local guild_capacity_increase_capacity_reached = 5
local guild_capacity_increase_failed = 6

function on_click(me, npc)
	local npc = 2010007
	local selected = me:dialog_list(npc,
		'길드를 만들고 싶은가? 혹은 길드 관련 업무를 위해서 찾아왔는가? 원하는 것을 말해보게.',
		{
			'길드를 만들고 싶습니다.',
			'길드를 해체합니다.',
			'길드 최대인원을 늘리고 싶습니다. (최대 100명)',
			'길드 최대인원을 늘리고 싶습니다. (최대 200명)',
		})
	if selected == nil then
		return
	end

	if selected == 0 then
		if me:guild() then
			me:dialog(npc, '흐음.. 이미 길드에 가입되어 있는 것 같은데?')
			return
		end
		if not me:dialog_yes_no(npc, '길드 제작 수수료는 #b1,500,000 메소#k라네, 정말 만들어 보고 싶은가?') then
			return
		end
		local result = me:generic_guild_message(1)
		if result == guild_create_ok then
			me:dialog(npc, '길드가 성공적으로 생성되었습니다.')
		elseif result == guild_create_not_allowed then
			me:dialog(npc, '길드를 만들 수 없습니다.')
		elseif result == guild_create_insufficient_meso then
			me:dialog(npc, '길드를 생성할 메소가 부족합니다.')
		elseif result == guild_create_already_in_guild then
			me:dialog(npc, '흐음.. 이미 길드에 가입되어 있는 것 같은데?')
		end
	elseif selected == 1 then
		local g = me:guild()
		if not g or g:rank(me) ~= 1 then
			me:dialog(npc, '길드장만이 길드를 해체할 수 있다네.')
			return
		end
		if not me:dialog_yes_no(npc, '길드를 해체하고 싶은가....? 지금 해체하게 된다면 절대 되돌릴 수 없다네.. 또, 모아뒀던 GP는 모두 사라지게 된다네. 길드 해체는 신중하게 선택하도록 하게나. 다시 한번 생각해보기 바라네. 정말 길드를 해체하고 싶은가?') then
			return
		end
		local result = g:disband(me)
		if result == guild_disband_not_master then
			me:dialog(npc, '길드장만이 길드를 해체할 수 있다네.')
		elseif result == guild_disband_not_in_guild then
			me:dialog(npc, '길드장만이 길드를 해체할 수 있다네.')
		elseif result == guild_disband_failed then
			me:dialog(npc, '길드 해체에 실패했습니다.')
		end
	elseif selected == 2 then
		local g = me:guild()
		if not g or g:rank(me) ~= 1 then
			me:dialog(npc, '길드장만이 길드 인원을 늘릴 수 있다네.')
			return
		end
		local extended_cap = false
		if not me:dialog_yes_no(npc, '길드 최대 인원 추가 비용은 #b50만#k 메소 라네. 지금 추가하면 최대 인원이 5명 만큼 더 늘어날걸세. 정말 최대 인원을 늘려보고 싶은가?') then
			return
		end
		local result = g:inc_capacity(me, extended_cap)
		if result == guild_capacity_increase_ok then
			-- me:dialog(npc, '길드 최대 인원이 증가했습니다.')
		elseif result == guild_capacity_increase_insufficient_meso then
			me:dialog(npc, '자네.. 메소는 충분히 갖고 있는건가?')
		elseif result == guild_capacity_increase_not_master then
			me:dialog(npc, '길드장만이 길드 인원을 늘릴 수 있다네.')
		elseif result == guild_capacity_increase_not_in_guild then
			me:dialog(npc, '길드장만이 길드 인원을 늘릴 수 있다네.')
		elseif result == guild_capacity_increase_capacity_reached then
			me:dialog(npc, '이미 길드 최대 인원 제한인 100 명이 된 것 같군.')
		else
			me:dialog(npc, '길드 최대 인원 증가에 실패했습니다.')
		end
	elseif selected == 3 then
		local g = me:guild()
		if not g or g:rank(me) ~= 1 then
			me:dialog(npc, '길드장만이 길드 인원을 늘릴 수 있다네.')
			return
		end
		local extended_cap = true
		if not me:dialog_yes_no(npc, '길드 최대 인원 추가 비용은 #b2,000#k 길드포인트 라네. 지금 추가하면 최대 인원이 5명 만큼 더 늘어날걸세. 정말 최대 인원을 늘려보고 싶은가?') then
			return
		end
		local result = g:inc_capacity(me, extended_cap)
		if result == guild_capacity_increase_ok then
			-- me:dialog(npc, '길드 최대 인원이 증가했습니다.')
		elseif result == guild_capacity_increase_not_master then
			me:dialog(npc, '길드장만이 길드 인원을 늘릴 수 있다네.')
		elseif result == guild_capacity_increase_not_in_guild then
			me:dialog(npc, '길드장만이 길드 인원을 늘릴 수 있다네.')
		elseif result == guild_capacity_increase_insufficient_gp then
			me:dialog(npc, '길드 포인트가 충분하지 않다네.')
		elseif result == guild_capacity_increase_capacity_reached then
			me:dialog(npc, '이미 길드 최대 인원 제한인 200 명이 된 것 같군.')
		else
			me:dialog(npc, '길드 최대 인원 증가에 실패했습니다.')
		end
	end
end
