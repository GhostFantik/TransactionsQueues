CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

create table users
(
    id      uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name    varchar(254) NOT NULL,
    balance integer      NOT NULL
        CONSTRAINT greater_or_equal_0 CHECK ( balance >= 0 ) DEFAULT 0
);