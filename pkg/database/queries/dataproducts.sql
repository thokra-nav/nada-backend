-- name: GetDataproduct :many
SELECT *
FROM dataproduct_complete_view
WHERE dataproduct_id = @id;

-- name: GetMinimalDataproductsByIDs :many
SELECT *
FROM dataproducts
WHERE id = ANY (@ids::uuid[])
ORDER BY last_modified DESC;

-- name: GetDataproducts :many
SELECT *
FROM dataproduct_complete_view
ORDER BY dp_last_modified DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetDataproductsByIDs :many
SELECT *
FROM dataproduct_complete_view
WHERE id = ANY (@ids::uuid[])
ORDER BY dp_last_modified DESC;

-- name: GetDataproductsByGroups :many
SELECT *
FROM dataproduct_complete_view
WHERE "dp_group" = ANY (@groups::text[])
ORDER BY dp_last_modified DESC;

-- name: GetDataproductsByProductArea :many
SELECT *
FROM dataproduct_complete_view
WHERE team_id = ANY(@team_id::text[])
ORDER BY dp_created DESC;

-- name: GetDataproductsByTeam :many
SELECT *
FROM dataproduct_complete_view
WHERE team_id = @team_id
ORDER BY dp_created DESC;

-- name: DeleteDataproduct :exec
DELETE
FROM dataproducts
WHERE id = @id;

-- name: CreateDataproduct :one
INSERT INTO dataproducts ("name",
                          "description",
                          "group",
                          "teamkatalogen_url",
                          "slug",
                          "team_contact",
                          "team_id")
VALUES (@name,
        @description,
        @owner_group,
        @owner_teamkatalogen_url,
        @slug,
        @team_contact,
        @team_id)
RETURNING *;

-- name: UpdateDataproduct :one
UPDATE dataproducts
SET "name"              = @name,
    "description"       = @description,
    "slug"              = @slug,
    "teamkatalogen_url" = @owner_teamkatalogen_url,
    "team_contact"      = @team_contact,
    "team_id"           = @team_id
WHERE id = @id
RETURNING *;


-- name: DataproductKeywords :many
SELECT keyword::text, count(1) as "count"
FROM (
	SELECT unnest(ds.keywords) as keyword
	FROM dataproducts dp
    INNER JOIN datasets ds ON ds.dataproduct_id = dp.id
) keywords
WHERE true
AND CASE WHEN coalesce(TRIM(@keyword), '') = '' THEN true ELSE keyword ILIKE @keyword::text || '%' END
GROUP BY keyword
ORDER BY keywords."count" DESC
LIMIT 15;

-- name: DataproductGroupStats :many
SELECT "group",
       count(1) as "count"
FROM "dataproducts"
GROUP BY "group"
ORDER BY "count" DESC
LIMIT @lim OFFSET @offs;
