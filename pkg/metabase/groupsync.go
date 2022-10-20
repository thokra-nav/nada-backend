package metabase

import (
	"context"
	"errors"
	"net/http"
)

type MetabaseSetting struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type GroupMappingOperation string

const (
	GroupMappingOperationAdd    GroupMappingOperation = "add"
	GroupMappingOperationRemove GroupMappingOperation = "remove"
)

func (c *Client) UpdateGroupMapping(ctx context.Context, azureGroupID string, mbPermissionGroupID int, operation GroupMappingOperation) error {
	current, err := c.getGroupMappings(ctx)
	if err != nil {
		return err
	}

	var updated map[string]interface{}
	switch operation {
	case GroupMappingOperationAdd:
		updated = addGroupMapping(current, azureGroupID, mbPermissionGroupID)
	case GroupMappingOperationRemove:
		updated = removeGroupMapping(current, azureGroupID, mbPermissionGroupID)
	default:
		return errors.New("invalid group mapping operation")
	}

	payload := map[string]map[string]interface{}{"saml-group-mappings": updated}
	if err := c.request(ctx, http.MethodPut, "/setting", payload, nil); err != nil {
		return err
	}

	return nil
}

func addGroupMapping(mappings map[string]interface{}, azureGroupID string, mbPermissionGroupID int) map[string]interface{} {
	for aGroup, pGroups := range mappings {
		if aGroup == azureGroupID {
			mappings[aGroup] = addGroup(pGroups.([]interface{}), mbPermissionGroupID)
			return mappings
		}
	}
	mappings[azureGroupID] = []int{mbPermissionGroupID}

	return mappings
}

func addGroup(groups []interface{}, group int) []interface{} {
	for _, g := range groups {
		if int(g.(float64)) == group {
			return groups
		}
	}
	groups = append(groups, group)

	return groups
}

func removeGroupMapping(mappings map[string]interface{}, azureGroupID string, mbPermissionGroupID int) map[string]interface{} {
	for aGroup, pGroups := range mappings {
		if aGroup == azureGroupID {
			mappings[aGroup] = removeGroup(pGroups.([]interface{}), mbPermissionGroupID)
			return mappings
		}
	}

	return mappings
}

func removeGroup(groups []interface{}, group int) []interface{} {
	for idx, g := range groups {
		if int(g.(float64)) == group {
			return append(groups[:idx], groups[idx+1:]...)
		}
	}

	return groups
}

func (c *Client) getGroupMappings(ctx context.Context) (map[string]interface{}, error) {
	settings := []*MetabaseSetting{}
	if err := c.request(ctx, http.MethodGet, "/setting", nil, &settings); err != nil {
		return nil, err
	}

	return getSAMLMappingFromSettings(settings)
}

func getSAMLMappingFromSettings(settings []*MetabaseSetting) (map[string]interface{}, error) {
	for _, s := range settings {
		if s.Key == "saml-group-mappings" {
			return s.Value.(map[string]interface{}), nil
		}
	}
	return nil, errors.New("saml group mappings not found in metabase settings")
}
