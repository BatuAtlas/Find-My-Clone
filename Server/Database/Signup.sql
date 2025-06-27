-- Table: public.Signup

-- DROP TABLE IF EXISTS public."Signup";

CREATE TABLE IF NOT EXISTS public."Signup"
(
    mail character varying(320) COLLATE pg_catalog."default" NOT NULL,
    pass character varying(64) COLLATE pg_catalog."default" NOT NULL,
    nickname character varying(64) COLLATE pg_catalog."default" NOT NULL,
    token character varying(64) COLLATE pg_catalog."default" NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."Signup"
    OWNER to postgres;