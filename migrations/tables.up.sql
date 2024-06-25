create table users
(
    id serial primary key ,
    login varchar(225) unique ,
    key varchar(225),
    name varchar(225),
    description text default 'O вас ...',
    avatar varchar(225) default 'none'
)