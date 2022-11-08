create table users_balance (
    id_user    int not null,
    balance    float(32) not null,
    constraint id_user unique (id_user)
);
grant all on users_balance to public;

create table reservation_money (
    id_reservation serial primary key,
    id_user        int not null,
    id_service     int not null,
    id_order       int not null,
    reserved_money float(32) not null,
    constraint unique_user_order unique (id_user, id_order)
);
grant all on reservation_money to public;
grant usage, select on sequence reservation_money_id_reservation_seq to public;

create table company_revenue (
    id_user      int not null,
    id_service   int not null,
    id_order     int not null,
    revenue      float(32) not null,
    date_create  timestamp default current_timestamp
);
grant all on company_revenue to public;