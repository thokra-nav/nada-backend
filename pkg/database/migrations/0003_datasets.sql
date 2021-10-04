-- +goose Up
CREATE TABLE datasets (
    "id" uuid DEFAULT uuid_generate_v4(),
		"dataproduct_id" uuid NOT NULL,
    "name" TEXT NOT NULL,
    "description" TEXT,
    "pii" BOOLEAN NOT NULL,
    "created" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "last_modified" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "project_id" TEXT NOT NULL,
    "dataset" TEXT NOT NULL,
    "table_name" TEXT NOT NULL,
    "type" TEXT NOT NULL,
    PRIMARY KEY(id),
		CONSTRAINT fk_dataproduct
			FOREIGN KEY(dataproduct_id)
				REFERENCES dataproducts(id) 
);

CREATE TRIGGER datasets_set_modified
BEFORE UPDATE ON datasets
FOR EACH ROW
EXECUTE PROCEDURE update_modified_timestamp();

-- +goose Down
DROP TABLE datasets;
DROP TRIGGER datasets_set_modified;
