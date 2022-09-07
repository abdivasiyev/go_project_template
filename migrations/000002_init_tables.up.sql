CREATE TABLE IF NOT EXISTS "company"
(
    "id"         uuid PRIMARY KEY NOT NULL,
    "name"       varchar          NOT NULL,
    "created_at" timestamp        NOT NULL,
    "updated_at" timestamp,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "user"
(
    "id"            uuid PRIMARY KEY NOT NULL,
    "company_id"    uuid references company (id),
    image_id        uuid references file (id),
    "username"      varchar UNIQUE   NOT NULL,
    "password_hash" varchar          NOT NULL,
    "first_name"    varchar,
    "last_name"     varchar,
    "phone"         varchar,
    "created_at"    timestamp        NOT NULL,
    "updated_at"    timestamp,
    "deleted_at"    timestamp
);

CREATE TABLE IF NOT EXISTS "role"
(
    "id"          uuid PRIMARY KEY NOT NULL,
    "alias"       varchar UNIQUE   NOT NULL,
    "name"        varchar          NOT NULL,
    "description" text,
    "is_basic"    boolean          not null default false,
    "created_at"  timestamp        NOT NULL,
    "updated_at"  timestamp,
    "deleted_at"  timestamp
);

create table if not exists permission
(
    id                uuid primary key not null default uuid_generate_v4(),
    alias             varchar          not null,
    sequence          varchar          not null,
    name              varchar          not null,

    path              varchar          not null default '/',
    method            varchar          not null default 'GET',
    query_param       varchar,
    query_param_value varchar,

    allow_all         boolean                   default false,

    created_at        timestamp        not null default current_timestamp,
    updated_at        timestamp,
    deleted_at        timestamp
);

CREATE TABLE IF NOT EXISTS "permission_module"
(
    "id"         uuid PRIMARY KEY NOT NULL,
    "alias"      varchar UNIQUE   NOT NULL,
    "name"       varchar          NOT NULL,
    "sequence"   integer          NOT NULL,
    "created_at" timestamp        NOT NULL,
    "updated_at" timestamp,
    "deleted_at" timestamp
);

CREATE TABLE IF NOT EXISTS "permission_group"
(
    "id"         uuid PRIMARY KEY                       NOT NULL,
    "module_id"  uuid references permission_module (id) NOT NULL,
    "alias"      varchar UNIQUE                         NOT NULL,
    "name"       varchar                                NOT NULL,
    "sequence"   integer                                NOT NULL,
    "created_at" timestamp                              NOT NULL,
    "updated_at" timestamp,
    "deleted_at" timestamp
);


create table if not exists permission_group_relation
(
    permission_id uuid not null references permission (id),
    group_id      uuid not null references permission_group (id)
);

create index if not exists idx_permission_path_and_method on permission (path, method);

CREATE TABLE IF NOT EXISTS "role_permission"
(
    "role_id"       uuid NOT NULL references "role" (id),
    "permission_id" uuid NOT NULL references permission (id)
);

CREATE TABLE IF NOT EXISTS "user_role"
(
    "user_id" uuid NOT NULL references "user" (id),
    "role_id" uuid NOT NULL references "role" (id)
);

insert into public.permission_module (id, alias, name, sequence, created_at, updated_at, deleted_at)
values ('e4749b62-3230-4ca4-ba27-ec4bcec43876', 'settings', 'Settings', 1, '2022-05-17 21:42:25.602887', null, null)
on conflict(id) do nothing;

insert into "user" (id, company_id, username, password_hash, first_name, last_name, created_at)
values (uuid_generate_v4(), null, 'admin',
        '$argon2id$v=19$m=65536,t=3,p=2$e+jCgaRufnSUCDGgdY48ew$fy0sTGeg8391Dg2AqWpFc5tpo6+KAta2pmKftahWxxg', 'Admin',
        'Admin', now())
on conflict(username) do nothing;
-- admin - qwerty123

CREATE TABLE IF NOT EXISTS file
(
    id         uuid primary key not null,
    name       varchar          not null,
    url        varchar          not null,
    created_at timestamp        not null,
    updated_at timestamp,
    deleted_at timestamp
);

create table if not exists app_version
(
    id           bigserial primary key not null,
    version      varchar               not null,
    title        varchar               not null default '',
    description  varchar               not null default '',
    force_update boolean               not null default false,

    created_at   timestamp             not null default current_timestamp
);

insert into app_version (id, version, title, description)
values (1, '1.0.0.0', 'Version 1.0.0', 'New Release')
on conflict (id) do nothing;