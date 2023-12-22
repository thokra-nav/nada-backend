-- name: MapDataset :exec
INSERT INTO third_party_mappings (
    "dataset_id",
    "services"
) VALUES (
    @dataset_id,
    @services
) ON CONFLICT ("dataset_id") DO UPDATE SET
    "services" = EXCLUDED.services;

-- name: GetDatasetMappings :one
SELECT *
FROM third_party_mappings
WHERE "dataset_id" = @dataset_id;

-- name: GetDatasetsByMapping :many
SELECT * FROM dataproduct_complete_view
WHERE @service::TEXT = ANY("mapping_services")
LIMIT @lim OFFSET @offs;
