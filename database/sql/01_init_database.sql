CREATE USER org_bot_user WITH PASSWORD 'orgBotPassword';

CREATE DATABASE org_bot
    WITH OWNER = org_bot_user;

GRANT pg_read_all_data TO org_bot_user;
GRANT pg_write_all_data TO org_bot_user;
