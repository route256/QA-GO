create type locale as ENUM ('RU', 'BY', 'KZ');
create table users
(
    id         serial             not null primary key,
    name       varchar(40)        not null,
    birthday   timestamptz,
    is_deleted bool default false not null,
    bio        text,
    locale     locale
);

insert into users (name, birthday, bio, locale)
values ('Dima', '1991-07-09 00:00:01+03', 'SDET TL from Ozon', 'RU');

select *
from users;

insert into users (name, birthday, bio, locale)
values ('Дункан Маклауд', '1592-03-02 04:00:01+03', 'Горец', 'BY'),
       ('Валентин', '2000-05-01 07:00:01+03', 'Священик', 'KZ');

delete
from users
where id = 3;

insert into users (name, birthday, bio, locale)
values ('Евгений', '1953-02-02 04:00:01+03', 'Водитель', 'RU'),
       ('Анастасия', '2001-05-01 07:00:01+03', 'Тренер', 'KZ');

update users
set is_deleted = true
where id = 4;

select count(true)
from users
where name = 'Dima';

select distinct locale
from users;

select *
from users
where name like '%д%';

create table items
(
    id          serial       not null primary key,
    title       varchar(256) not null,
    description text,
    price       numeric(15, 2)
        constraint positive_price check ( price > 0 ),
    stock       integer      not null
        constraint non_negative_stock check ( stock >= 0 )
);

insert into items (title, description, price, stock)
values ('Война и мир', 'Мировой бестселлер', 799, 5),
       ('Пряник', 'Тульский', 59.5, 100),
       ('Мазь "Звёздночка"', 'Лечит все белзни', 99, 1000);

select *
from items;

create table category
(
    id   serial not null primary key,
    name varchar(100)
);

insert into category (name)
values ('Книги'),
       ('Смартфоны'),
       ('Продукты'),
       ('Лекарства');

select *
from category;

alter table items
    add column category_id integer,
    add constraint fk_category_id
        foreign key (category_id)
            references category (id);

select *
from items;

update items
set category_id = 4
where id = 3;

update items
set title = 'Мазь "Звёздночка"'
where id = 3;

insert into items (title, description, price, stock, category_id)
values ('Iphone 14', 'Самый лучший смартфон', 99999.99, 10, 2),
       ('Iphone 14 Pro', 'Ещё более лучший смартфон', 109999.99, 10, 2),
       ('Iphone 14 Pro MAX"', 'Наилучший смартфон', 129999.99, 10, 2);

select AVG(price)
from items
where category_id = 2;

select AVG(price)
from items
where category_id = 2;

select MAX(price)
from items;

select MIN(price)
from items;

SELECT SUM(stock)
from items
where category_id = 2;

select AVG(price), category_id
from items
group by category_id;

select AVG(price), category_id
from items
group by category_id
having AVG(price) > 100;

SELECT SUM(price * stock)
from items;

select id, title, price
from items
order by price;

select id, title, price
from items
order by price desc;

select id, title, price
from items
order by price desc
limit 1;

select id, title, price
from items
where price = (select MAX(price)
               from items);


create table carts
(
    id      serial primary key,
    user_id integer references users (id)
);

select *
from users;

insert into carts (user_id)
values (1),
       (2),
       (5)

create table items_in_cart
(
    id      serial primary key,
    cart_id integer references carts (id),
    item_id integer references items (id),
    amount  integer not null
        constraint positive_amount check ( amount > 0)
);

select *
from items;

insert into items_in_cart (cart_id, item_id, amount)
values (1, 1, 1),
       (3, 2, 10),
       (3, 5, 1);

select *
from items_in_cart;

select ic.cart_id, i.title, ic.amount
from items_in_cart ic
         join items i on ic.item_id = i.id;

select u.name, i.title, ic.amount
from items_in_cart ic
         join items i on ic.item_id = i.id
         join carts c on c.id = ic.cart_id
         join users u on c.user_id = u.id;

select i.title
from items_in_cart ic
         right join items i on ic.item_id = i.id
where amount is null;