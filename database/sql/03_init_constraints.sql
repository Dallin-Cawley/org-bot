\connect org_bot

ALTER TABLE org_bot_schema.join_vc
    ADD CONSTRAINT fk_guild_id
        FOREIGN KEY (guild_id) REFERENCES org_bot_schema.guild(guild_id);