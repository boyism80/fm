local M = {}

M.REPLIES = {
	resp.dialog,
	resp.dialog_yes_no,
	resp.dialog_input,
	resp.dialog_list,
	resp.dialog_style,
	resp.dialog_accept,
	resp.warp,
	resp.update_stats,
	resp.open_npc_shop,
}

function M.replied(p, name)
	return name ~= resp.update_stats or p.unlock_action
end

function M.is_dialog(name)
	return name ~= resp.warp and name ~= resp.update_stats and name ~= resp.open_npc_shop
end

function M.queues(count, entries, check)
	local queues = {}
	for i = 0, count - 1 do
		queues[#queues + 1] = function(ctx)
			local bot = ctx:bot(i)
			for j = i + 1, #entries, count do
				check(ctx, bot, entries[j])
			end
			return true
		end
	end
	return queues
end

return M
