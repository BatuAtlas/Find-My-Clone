CREATE TABLE public."GPS"
(
    "user" bigint NOT NULL,
    lat double precision NOT NULL,
    lon double precision NOT NULL,
    elevation integer,
    "timestamp" timestamp with time zone NOT NULL,
    PRIMARY KEY ("user"),
    FOREIGN KEY ("user")
        REFERENCES public."User" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE IF EXISTS public."GPS"
    OWNER to postgres;