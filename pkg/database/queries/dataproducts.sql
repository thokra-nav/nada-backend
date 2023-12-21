-- name: GetDataproduct :one
SELECT *
FROM dataproducts
WHERE id = @id;

-- name: GetDataproducts :many
SELECT *
FROM dataproducts
ORDER BY last_modified DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetDataproductsByIDs :many
SELECT *
FROM dataproducts
WHERE id = ANY (@ids::uuid[])
ORDER BY last_modified DESC;

-- name: GetDataproductsByGroups :many
SELECT *
FROM dataproducts
WHERE "group" = ANY (@groups::text[])
ORDER BY last_modified DESC;

-- name: GetDataproductsByProductArea :many
SELECT *
FROM dataproducts
WHERE team_id = ANY(@team_id::text[])
ORDER BY created DESC;

-- name: GetDataproductsByTeam :many
SELECT *
FROM dataproducts
WHERE team_id = @team_id
ORDER BY created DESC;

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

-- name: GetDataproductComplete :many
SELECT 
dsrc.id AS dsrc_id,  
dsrc.created as dsrc_created,
dsrc.last_modified as dsrc_last_modified,
dsrc.expires as dsrc_expires,
dsrc.description as dsrc_description,
dsrc.missing_since as dsrc_missing_since,
dsrc.pii_tags as pii_tags,
dsrc.project_id as project_id,
dsrc.dataset as dataset,
dsrc.table_name as table_name,
dsrc.table_type as table_type,
dsrc.pseudo_columns as pseudo_columns,
dsrc.dataset_id as dsrc_dataset_id,
dsrc.schema as dsrc_schema,
dpds.*,
dm.services,
da.id as da_id,
da.subject as da_subject,
da.granter as da_granter,
da.expires as da_expires,
da.created as da_created,
da.revoked as da_revoked,
da.access_request_id as access_request_id,
mm.database_id as mm_database_id
FROM 
(
	SELECT 
    ds.id AS ds_id, 
    ds.name as ds_name, 
    ds.description as ds_description,
    ds.created as ds_created,
    ds.last_modified as ds_last_modified,
    ds.slug as ds_slug,
    ds.keywords as keywords,
    rdp.* 
    FROM 
	(
		(SELECT * FROM dataproducts dp WHERE dp.id= @id) rdp 
		LEFT JOIN datasets ds ON ds.dataproduct_id = rdp.id
	)
) dpds 
LEFT JOIN 
    (SELECT * FROM datasource_bigquery WHERE is_reference = false) dsrc
ON dpds.ds_id = dsrc.dataset_id 
LEFT JOIN third_party_mappings dm ON dpds.ds_id = dm.dataset_id
LEFT JOIN dataset_access da ON dpds.ds_id = da.dataset_id
LEFT JOIN metabase_metadata mm ON mm.dataset_id = dpds.ds_id AND mm.deleted_at IS NULL;
