\connect org_bot
CREATE SCHEMA IF NOT EXISTS org_bot_schema;

CREATE TABLE IF NOT EXISTS org_bot_schema.guild
(
    guild_id    text PRIMARY KEY,
    bot_role_id text,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

/* join_vc is the table containing the join to create channel ids */
CREATE TABLE IF NOT EXISTS org_bot_schema.join_vc
(
    join_vc_id  text PRIMARY KEY,
    guild_id    text,
    category_id text UNIQUE,
    created_at  TIMESTAMP DEFAULT NOW(),
    updated_at  TIMESTAMP
);

/* join_vc_child are channels spawned when a user joins a join_vc channel */
CREATE TABLE IF NOT EXISTS org_bot_schema.join_vc_child
(
    join_vc_child_id text PRIMARY KEY,
    guild_id         text,
    category_id      text,
    join_vc_id       text,
    created_at       TIMESTAMP DEFAULT NOW(),
    updated_at       TIMESTAMP
);