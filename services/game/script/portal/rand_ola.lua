local ola_ola = require("script/lib/event/ola_ola")

return {
    on_enter = function(me, portal)
        ola_ola.pass(me, portal)
    end
}
