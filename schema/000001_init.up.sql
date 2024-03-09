CREATE TABLE orders
(
    order_uid varchar(255) not null unique primary key, --PRIMARY KEY
    track_number varchar(255) not null unique, --PRIMARY KEY
    entry varchar(255) not null,
    locale varchar(2) not null,
    internal_signature varchar(255),
    customer_id varchar(255), --MAYBE PRIMARY KEY
    delivery_service varchar(255),
    shardkey varchar(255),
    sm_id int, --MAYBE WRONG
    date_created timestamp with time zone,
    oof_shard varchar(255)
);

CREATE TABLE deliveries
(
    order_uid varchar(255) references orders(order_uid) on delete cascade not null,
    name varchar(255) not null,
    phone varchar(255) not null,
    zip varchar(255) not null,
    city varchar(255) not null,
    address varchar(255) not null,
    region varchar(255) not null,
    email varchar(255) not null
);

CREATE TABLE payments
(
    transaction varchar(255) references orders(order_uid) on delete cascade not null,
    request_id varchar(255),
    currency varchar(255) not null,
    provider varchar(255) not null,
    amount int not null,
    payment_dt int not null,
    bank varchar(255) not null,
    delivery_cost int not null,
    goods_total int not null,
    custom_fee int not null
);

CREATE TABLE items
(
    order_uid varchar(255) references orders(order_uid) on delete cascade not null,
    chrt_id int not null unique primary key,
    track_number varchar(255) not null,
    price int not null,
    rid varchar(255) not null,
    name varchar(255) not null,
    sale int not null,
    size varchar(255) not null,
    total_price int not null,
    nm_id int not null unique,
    brand varchar(255) not null,
    status int not null
);