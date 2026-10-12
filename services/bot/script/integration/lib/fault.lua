local pq = require("script/integration/lib/party_quest")

local M = {}

function M.rpc(bot, method, mode)
	local text = string.format("/RPC장애 %s %s", method, mode)
	return pq.command(bot, text, string.format("RPC 장애(%s)", method)) ~= false
end

function M.rpc_delay(bot, method, ms)
	local text = string.format("/RPC지연 %s %d", method, ms)
	return pq.command(bot, text, string.format("RPC 지연(%s)을 %dms로", method, ms)) ~= false
end

return M
