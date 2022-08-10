-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.products (
    id bigserial primary key,
    name varchar(255) not null,
    price bigint not null CONSTRAINT positive_product_price CHECK (price >= 0),
    quantity bigint not null CONSTRAINT positive_product_quantity CHECK (quantity >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.products;
-- +goose StatementEnd
