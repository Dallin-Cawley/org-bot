\connect org_bot

/************************************
 * CREATE MODEL
 ***********************************/

CREATE OR REPLACE FUNCTION org_bot_schema.new_guild(_guild_id text,
                                                    _bot_role_id text)

    RETURNS TABLE (
                      guild_id text,
                      bot_role_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
    INSERT INTO org_bot_schema.guild(guild_id, bot_role_id)
        VALUES(_guild_id, _bot_role_id)
        ON CONFLICT DO NOTHING
    RETURNING *;
$BODY$
LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.new_join_vc(_join_vc_id text,
                                                      _guild_id text,
                                                      _category_id text)

    RETURNS TABLE (
                      join_vc_id text,
                      guild_id text,
                      category_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
    INSERT INTO org_bot_schema.join_vc(join_vc_id, guild_id, category_id)
           VALUES(_join_vc_id, _guild_id, _category_id)
    RETURNING *;
$BODY$
LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.new_join_vc_child(_join_vc_child_id text,
                                                            _guild_id text,
                                                            _category_id text,
                                                            _join_vc_id text)

    RETURNS TABLE (
                      join_vc_child_id text,
                      guild_id text,
                      category_id text,
                      join_vc_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
INSERT INTO org_bot_schema.join_vc_child(join_vc_child_id, guild_id, category_id, join_vc_id)
    VALUES(_join_vc_child_id, _guild_id, _category_id, _join_vc_id)
RETURNING *;
$BODY$
    LANGUAGE SQL;

/************************************
 * READ MODEL
 ***********************************/

CREATE OR REPLACE FUNCTION org_bot_schema.read_guild(_guild_id text)

    RETURNS TABLE (
                      guild_id text,
                      bot_role_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
SELECT * FROM org_bot_schema.guild
    WHERE guild_id = _guild_id;
$BODY$
LANGUAGE SQL;


CREATE OR REPLACE FUNCTION org_bot_schema.read_join_vc(_join_vc_id text)

    RETURNS TABLE (
                      join_vc_id text,
                      guild_id text,
                      category_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
SELECT * FROM org_bot_schema.join_vc
    WHERE join_vc_id = _join_vc_id;
$BODY$
LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.read_join_vc_child(_join_vc_child_id text)

    RETURNS TABLE (
                      join_vc_child_id text,
                      guild_id text,
                      category_id text,
                      join_vc_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
SELECT * FROM org_bot_schema.join_vc_child
    WHERE join_vc_child_id = _join_vc_child_id;
$BODY$
LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.read_join_vc_in_category(_category_id text)

    RETURNS TABLE (
                      join_vc_id text,
                      guild_id text,
                      category_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
SELECT * FROM org_bot_schema.join_vc
WHERE category_id = _category_id;
$BODY$
    LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.read_guild_join_vc(_guild_id text)

    RETURNS TABLE (
                      join_vc_id text,
                      guild_id text,
                      category_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
SELECT * FROM org_bot_schema.join_vc
    WHERE guild_id = _guild_id;
$BODY$
LANGUAGE SQL;

/************************************
 * DELETE MODEL
 ***********************************/

CREATE OR REPLACE FUNCTION org_bot_schema.delete_guild(_guild_id text)

    RETURNS TABLE (
                      guild_id text,
                      bot_role_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
DELETE FROM org_bot_schema.guild
    WHERE guild_id = _guild_id
RETURNING *;
$BODY$
    LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.delete_join_vc(_join_vc_id text)

    RETURNS TABLE (
                      join_vc_id text,
                      guild_id text,
                      category_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
DELETE FROM org_bot_schema.join_vc
    WHERE join_vc_id = _join_vc_id
RETURNING *;
$BODY$
LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.delete_join_vc_child(_join_vc_child_id text)

    RETURNS TABLE (
                      join_vc_child_id text,
                      guild_id text,
                      category_id text,
                      join_vc_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
DELETE FROM org_bot_schema.join_vc_child
WHERE join_vc_child_id = _join_vc_child_id
RETURNING *;
$BODY$
    LANGUAGE SQL;

CREATE OR REPLACE FUNCTION org_bot_schema.delete_guild_join_vc(_guild_id text)

    RETURNS TABLE (
                      join_vc_id text,
                      guild_id text,
                      category_id text,
                      created_at timestamp,
                      updated_at timestamp
                  ) AS
$BODY$
DELETE FROM org_bot_schema.join_vc
    WHERE guild_id = _guild_id
RETURNING *;
$BODY$
LANGUAGE SQL;