CREATE TABLE public."Authorization"
(
    token character varying(64) NOT NULL,
    "user" bigint NOT NULL, /*user id*/
    expires timestamp with time zone,
    PRIMARY KEY (token)
);

ALTER TABLE IF EXISTS public."Authorization"
    OWNER to postgres;