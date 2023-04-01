-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wdiet.ingredients
(
    ingredient_uuid uuid not null default gen_random_uuid()
        constraint ingredients_primary_key
            primary key,
    ingredient_name      varchar(64)     not null UNIQUE,
    category             varchar(64)     not null,
    days_until_exp       integer         not null,
    created_at           timestamp       not null default now(),
    updated_at           timestamp       not null default now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS wdiet;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wdiet.fridge_ingredients --change the way the data looks.
(
    user_uuid              uuid            not null
        constraint user_uuid_fk references wdiet.users,
    ingredient_uuid        uuid            not null
        constraint ingredient_uuid_fk references wdiet.ingredients,
    amount                 integer         not null,
    unit                   varchar(64)     not null,
    purchased_date         timestamp       not null, --expiry_date는 days until expire랑은 달라.. 그건 integer고 이건 date임. 근데 expiry_date도 하지마. ingredient의 days until exp를 이용해서 자동계산하려면 구매한 날짜만 알면됨. 그럼 유통기한은 자동 계산되니까.
    expiration_date        timestamp       not null,
    created_at             timestamp       not null default now(),
    updated_at             timestamp       not null default now(),
    PRIMARY KEY (user_uuid, ingredient_uuid) --PRIMARY KEY doesn't need CREATE INDEX ON
);

-- CREATE INDEX ON wdiets.fridge_ingredients (user_uuid);
-- CREATE INDEX ON wdiets.fridge_ingredients (ingredient_uuid);

-- ALTER TABLE wdiet.fridge_ingredients ADD FOREIGN KEY user_uuid REFERENCES wdiet.users (user_uuid);
-- ALTER TABLE wdiet.fridge_ingredients ADD FOREIGN KEY ingredient_uuid REFERENCES wdiet.ingredients (ingredient_uuid);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wdiet.fridge_ingredients;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS wdiet.ingredients;
-- +goose StatementEnd