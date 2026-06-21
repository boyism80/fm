-- NPC name (String.wz/Npc.img.xml): 레아

local guild_emblem_dialog = 17

function on_start(me)
	local npc = 2010008
	local selected = me:dialog_list(npc,
		'저는 길드 마크 제작 업무를 맡고 있습니다. 길드 마크는 길드장만 변경할 수 있답니다. 원하는 것이 있으세요?',
		{
			'길드마크 추가/변경',
		})
	if selected == nil or selected ~= 0 then
		return
	end

	local g = me:guild()
	if not g or g:rank(me) ~= 1 then
		me:dialog(npc, '길드장만 길드 마크를 추가하거나 변경할 수 있답니다. 당신은 길드장이 아닌 것 같군요.')
		return
	end

	if not me:dialog_yes_no(npc, '길드마크 추가/변경 수수료는 #b500만 메소#k 입니다. 길드 마크를 제작해 보시고 싶으신가요?') then
		return
	end

	me:generic_guild_message(guild_emblem_dialog)
end
