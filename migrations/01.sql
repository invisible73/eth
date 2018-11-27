create table transactions (
    id SERIAL,
    "date" timestamp  DEFAULT now(),
    "from" char(42) not null,
    "to" char(42) not null,
    "amount" numeric(50,0),
    "hash" char(66) not null,
    "block_hash" char(66),
    "block_number" integer,
    confirmations integer not null default 0,
    invalidated bool not null DEFAULT false,
    sended bool not null DEFAULT false,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

create index transactions_confirmations_idx on transactions(confirmations);