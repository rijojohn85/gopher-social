CREATE TABLE if not exists user_invitations (
    token bytea primary key,
    user_id bigint not null
);