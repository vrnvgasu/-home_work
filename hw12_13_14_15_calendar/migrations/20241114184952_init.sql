-- +goose Up
-- +goose StatementBegin

create table event
(
    id          serial primary key,
    title       varchar(255) not null,
    start_at    timestamptz  not null,
    end_at      timestamptz  not null,
    description text,
    owner_id    bigint not null,
    send_before bigint
);
comment on column event.title is 'Заголовок';
comment on column event.start_at is 'Дата и время события';
comment on column event.end_at is 'Дата и время окончания';
comment on column event.description is 'Описание события';
comment on column event.owner_id is 'ID пользователя, владельца события';
comment on column event.send_before is 'За сколько времени высылать уведомление';

create index idx_books_description
    on event (start_at);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table if exists event;
-- +goose StatementEnd
