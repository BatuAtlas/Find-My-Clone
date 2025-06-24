CREATE TABLE IF NOT EXISTS public."Userinfo"
(
    "user" bigint NOT NULL,
    status character varying(50) COLLATE pg_catalog."default",
    "isCharging" boolean,
    battery bytea,
    event bytea,
    "lastUpdate" time with time zone,
    CONSTRAINT "Userinfo_pkey" PRIMARY KEY ("user"),
    CONSTRAINT "Userinfo_user_fkey" FOREIGN KEY ("user")
        REFERENCES public."User" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        NOT VALID
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Userinfo"
    OWNER to postgres;