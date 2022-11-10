create table users (
   user_id serial primary key,
   user_name varchar(30) not null,
   password  varchar(255) not null
);