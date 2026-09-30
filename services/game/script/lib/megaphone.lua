local JAIL_MAP = 180000002
local MAX_TEXT = 65

local M = {}

local function length(text)
	local _, count = string.gsub(text, "[^\128-\191]", "")
	return count
end

local function medal(me)
	local item = me:equipped(EquipmentPart.Medal)
	if item == nil then
		return ""
	end
	local name = item2name(item:wz():id())
	if name == nil then
		return ""
	end
	return "<" .. string.gsub(name, "의 훈장", "") .. "> "
end

function M.send(me, text, msg_type, scope, ear, opts)
	opts = opts or {}
	if opts.min_level ~= nil and me:level() < opts.min_level then
		me:message(string.format("레벨 %d 이상만 사용할 수 있습니다.", opts.min_level), Msg.PinkText)
		return false
	end
	if me:map():wz():id() == JAIL_MAP then
		me:message("이곳에서 사용할 수 없습니다.", Msg.PinkText)
		return false
	end
	local cooldown = opts.cooldown ~= nil and me:role() < ROLE.Admin
	if cooldown and now() - me:last_megaphone() < opts.cooldown then
		me:message(string.format("%d초에 한 번만 사용할 수 있습니다.", opts.cooldown), Msg.PinkText)
		return false
	end
	if get_megaphone_muted() then
		me:message("현재 확성기 사용 금지 상태입니다.", Msg.PinkText)
		return false
	end
	if text == nil or length(text) == 0 or length(text) > MAX_TEXT then
		return false
	end

	local ok = me:message(medal(me) .. me:name() .. " : " .. text, msg_type, scope, ear)
	if ok and cooldown then
		me:last_megaphone(now())
	end
	return ok
end

return M
