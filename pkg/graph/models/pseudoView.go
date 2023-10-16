package models

// NewPseudoView contains metadata for creating a new pseudonymised view
type NewPseudoView struct {
	// projectID is the GCP project ID of the target table.
	ProjectID string `json:"projectID"`
	// dataset is the name of the dataset of the target table.
	Dataset string `json:"dataset"`
	// table is the name of the target table
	Table string `json:"table"`
	// targetColumns is the columns to be pseudonymised.
	TargetColumns []string `json:"targetColumns,omitempty"`
}

// NewJoinableViews contains metadata for creating joinable views
type NewJoinableViews struct {
	// datasetIDs is the IDs of the dataset which connects to joinable views.
	DatasetIDs []string `json:"datasetIDs,omitempty"`
}
