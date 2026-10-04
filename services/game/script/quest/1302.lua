local carnival_rank = require("script/lib/carnival_rank")

return carnival_rank.hooks(1302, { "have" }, { "gvup", "vic", "lose", "draw" })
