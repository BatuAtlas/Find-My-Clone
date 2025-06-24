CREATE TABLE public."User"
(
    id bigserial NOT NULL,
    nickname character varying(45)[] NOT NULL,
    profilephoto text,
    friends bigint[],
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public."User"
    OWNER to postgres;