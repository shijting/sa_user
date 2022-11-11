-- 角色表
create table roles (
   id serial primary key,
   name varchar(30) not null,
   code varchar(30) not null,
   description varchar(500) null default ''
);

-- 用户角色表
create table user_roles (
    id serial primary key,
    user_id UUID not null,
    role_id INT not null
);

-- 权限表
create table permissions (
     id serial primary key,
     name varchar(30) not null,
     code varchar(30) not null,
     action varchar(255) not null,
     parent_id int not null DEFAULT 0,
     description varchar(500) null default ''
);

create table roles_permissions (
   id serial primary key,
   role_code varchar(50) not null,
   permission_code varchar(50) not null
);