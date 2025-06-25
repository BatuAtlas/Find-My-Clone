-- Table: public.User

-- DROP TABLE IF EXISTS public."User";

CREATE TABLE IF NOT EXISTS public."User"
(
    id bigint NOT NULL DEFAULT nextval('"User_id_seq"'::regclass),
    nickname character varying(64) COLLATE pg_catalog."default" NOT NULL,
    profilephoto text COLLATE pg_catalog."default",
    friends bigint[],
    CONSTRAINT "User_pkey" PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public."User"
    OWNER to postgres;