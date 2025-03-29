box.cfg{
    listen = 3301,
    memtx_memory = 128 * 1024 * 1024,
    log_level = 5
}

box.schema.user.grant('guest', 'read,write,execute', 'universe')

box.schema.space.create('polls', {
    if_not_exists = true,
    format = {
        {name = 'id', type = 'string'},
        {name = 'question', type = 'string'},
        {name = 'creator', type = 'string'},
        {name = 'channel_id', type = 'string'},
        {name = 'created_at', type = 'unsigned'},
        {name = 'status', type = 'string'},
        {name = 'options', type = 'map'}
    }
})

box.space.polls:create_index('primary', {
    type = 'hash',
    parts = {'id'},
    if_not_exists = true
})

box.schema.space.create('votes', {
    if_not_exists = true,
    format = {
        {name = 'poll_id', type = 'string'},
        {name = 'user_id', type = 'string'},
        {name = 'option_id', type = 'string'},
        {name = 'voted_at', type = 'unsigned'}
    }
})

box.space.votes:create_index('primary', {
    type = 'tree',
    parts = {'poll_id', 'user_id'},
    if_not_exists = true
})