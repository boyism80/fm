local create = require("script/lib/fourth_skill_instance")

return create({
    leave_on_party = true,
    map = 910500200,
    exit_map = 105090800,
    duration_ms = 1200000,
    max_machines = 1,
    -- Only one of the five stone shields holds Seiram's shield, and it moves each time.
    shuffle_reactors = true,
})
