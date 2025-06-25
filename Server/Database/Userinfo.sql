-- Table: public.Userinfo

-- DROP TABLE IF EXISTS public."Userinfo";

CREATE TABLE IF NOT EXISTS public."Userinfo"
(
    "user" bigint NOT NULL,
    status character varying(200) COLLATE pg_catalog."default",
    "isCharging" boolean,
    battery smallint,
    event smallint,
    "lastUpdate" timestamp with time zone NOT NULL,
    CONSTRAINT "Userinfo_pkey" PRIMARY KEY ("user"),
    CONSTRAINT "Userinfo_user_fkey" FOREIGN KEY ("user")
        REFERENCES public."User" (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT battery_range CHECK (battery >= 0 AND battery <= 100)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Userinfo"
    OWNER to postgres;