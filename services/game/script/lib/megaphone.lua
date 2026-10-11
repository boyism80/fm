local JAIL_MAP = 180000002
local MAX_TEXT_BYTES = 60

local M = {}

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

local function usable(me)
	if me:map():wz():id() == JAIL_MAP then
		me:message("이곳에서 사용할 수 없습니다.", Msg.PinkText)
		return false
	end
	if get_megaphone_muted() then
		me:message("현재 확성기 사용 금지 상태입니다.", Msg.PinkText)
		return false
	end
	return true
end

function M.send(me, text, msg_type, scope, ear)
	if not usable(me) then
		return false
	end
	if text == nil or not valid_text(text, MAX_TEXT_BYTES) then
		return false
	end
	return me:message(medal(me) .. me:name() .. " : " .. text, msg_type, scope, ear)
end

function M.send_lines(me, lines, ear)
	if not usable(me) then
		return false
	end
	if #lines < 1 or #lines > 3 then
		return false
	end
	local prefix = medal(me) .. me:name() .. " : "
	local prefixed = {}
	for _, line in ipairs(lines) do
		if not valid_text(line, MAX_TEXT_BYTES) then
			return false
		end
		table.insert(prefixed, prefix .. line)
	end
	return me:megaphone_lines(prefixed, ear)
end

function M.send_item(me, text, ear, inventory, slot)
	if not usable(me) then
		return false
	end
	if text == nil or not valid_text(text, MAX_TEXT_BYTES) then
		return false
	end
	return me:item_megaphone(medal(me) .. me:name() .. " : " .. text, ear, inventory, slot)
end

return M
