local M = {}

local WATER = 5180000
local SCROLL = 4031034
local NO_PET = "저는 생명의 시간이 다 되어 죽은 펫을 되살려 주고 있지만.. 현재 여행자님께선 되살릴 펫이 있는 것 같지 않군요.."

function M.talk(me, npc, quest_id, intro, materials)
	local sel = me:dialog_list(npc, intro, { "죽은 펫을 다시 살리고 싶어요." })
	if sel == nil then
		return
	end

	local pets = me:expired_pets()
	if #pets == 0 then
		me:dialog(npc, NO_PET)
		return
	end

	local q = me:quest(quest_id)
	if q == nil then
		return
	end
	if q:started() == false then
		if not me:dialog_yes_no(npc, "펫을 되살리고 싶으신건가요..? 좋아요. 그러기 위해선 준비물이 몇가지 필요합니다만..") then
			return
		end
		me:dialog(npc, materials)
		q:start(npc, true)
		return
	end
	if next(me:item(SCROLL)) == nil or next(me:item(WATER)) == nil then
		me:dialog(npc, materials)
		return
	end
	if not me:dialog_yes_no(npc, "준비물을 모두 모아오셨군요.. 좋아요. 이제 당신의 펫을 살릴 준비가 되었어요.") then
		return
	end

	local options = {}
	for _, pet in ipairs(pets) do
		table.insert(options, "#v" .. pet.id .. "# #t" .. pet.id .. "#")
	end
	local choice = me:dialog_list(npc, "살리고 싶은 펫을 선택해 주세요.", options)
	if choice == nil then
		return
	end
	if me:exchange({ item = { [WATER] = 1, [SCROLL] = 1 } }, nil) ~= ExchangeResult.OK then
		me:dialog(npc, materials)
		return
	end
	me:revive_pet(pets[choice].slot)
	q:force_complete(npc)
	me:dialog(npc, "펫이 되살아났군요. 좋아요.. 다시 태어난 당신의 친구와 행복한 시간을 보내길..")
end

return M
