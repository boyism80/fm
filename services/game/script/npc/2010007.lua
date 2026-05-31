local guild_create_ok = 0
local guild_create_not_allowed = 1
local guild_create_insufficient_meso = 2
local guild_create_already_in_guild = 3
local guild_create_failed = 4

local guild_disband_ok = 0
local guild_disband_not_in_guild = 1
local guild_disband_not_master = 2
local guild_disband_failed = 3

function on_start(me)
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
		if me:guild_id() then
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
		if not me:guild_id() or me:guild_rank() ~= 1 then
			me:dialog(npc, '길드장만이 길드를 해체할 수 있다네.')
			return
		end
		if not me:dialog_yes_no(npc, '길드를 해체하고 싶은가....? 지금 해체하게 된다면 절대 되돌릴 수 없다네.. 또, 모아뒀던 GP는 모두 사라지게 된다네. 길드 해체는 신중하게 선택하도록 하게나. 다시 한번 생각해보기 바라네. 정말 길드를 해체하고 싶은가?') then
			return
		end
		local result = disband_guild(me)
		if result == guild_disband_not_master then
			me:dialog(npc, '길드장만이 길드를 해체할 수 있다네.')
		elseif result == guild_disband_not_in_guild then
			me:dialog(npc, '길드장만이 길드를 해체할 수 있다네.')
		elseif result == guild_disband_failed then
			me:dialog(npc, '길드 해체에 실패했습니다.')
		end
	elseif selected == 2 then
		me:dialog(npc, '아직 지원하지 않습니다.')
	elseif selected == 3 then
		me:dialog(npc, '아직 지원하지 않습니다.')
	end
end
