WITH new_user AS (
  INSERT INTO "User" ("nickname")
  VALUES ('Selin')
  RETURNING "id"
),
userinfo_insert AS (
  INSERT INTO "Userinfo" ("user", "lastUpdate")
  SELECT "id", NOW()
  FROM new_user
),
auth_insert AS (
  INSERT INTO "Authorization" ("user", "token")
  SELECT "id", encode(gen_random_bytes(32), 'hex')
  FROM new_user
  RETURNING "token", "user"
)
SELECT "user" AS "id", "token" FROM auth_insert;
