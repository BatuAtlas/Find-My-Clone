SELECT 
  "user" as "id",
  CASE 
    WHEN password = $1 THEN true
    ELSE false
  END AS "match"
FROM "Usersettings"
WHERE "mail" = $2;