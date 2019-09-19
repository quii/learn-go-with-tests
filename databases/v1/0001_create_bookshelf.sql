
create table if not exists bookshelf.books
(
    id     serial  not null
        constraint books_pk
            primary key,
    title  varchar not null,
    author varchar not null
);

alter table bookshelf.books
    owner to postgres;

create unique index if not exists books_id_uindex
    on bookshelf.books (id);
