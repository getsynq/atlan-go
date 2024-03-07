package client

import (
	"atlan-go/atlan/model/assets"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	MaxRetries    = 5
	RetryInterval = time.Second * 5
)

// GlossaryClient defines the client for interacting with the model API.
type GlossaryClient struct {
	client *AtlanClient
}

// NewGlossaryClient creates a new GlossaryClient instance.
func NewGlossaryClient(ac *AtlanClient) *GlossaryClient {
	return &GlossaryClient{client: ac}
}

// GetGlossaryByGuid retrieves a glossary by its GUID.
func GetGlossaryByGuid(glossaryGuid string) (*assets.AtlasGlossary, error) {
	if DefaultAtlanClient == nil {
		return nil, fmt.Errorf("default AtlanClient not initialized")
	}

	api := &GET_ENTITY_BY_GUID
	api.Path += glossaryGuid

	response, err := DefaultAtlanClient.CallAPI(api, nil, nil)
	if err != nil {
		return nil, err
	}

	g, err := assets.FromJSON(response)
	if err != nil {
		return nil, err
	}
	fmt.Println("Glossary:", g)

	return g, nil
}

// GetGlossaryTermByGuid retrieves a glossary term by its GUID.
func GetGlossaryTermByGuid(glossaryGuid string) (*assets.GlossaryTerm, error) {
	if DefaultAtlanClient == nil {
		return nil, fmt.Errorf("default AtlanClient not initialized")
	}

	api := &GET_ENTITY_BY_GUID
	api.Path += glossaryGuid

	response, err := DefaultAtlanClient.CallAPI(api, nil, nil)
	if err != nil {
		return nil, err
	}

	gt, err := assets.FromJSONTerm(response)
	if err != nil {
		return nil, err
	}
	fmt.Println("GlossaryTerm:", gt)

	return gt, nil
}

// Creator is used to create a new glossary asset in memory.
func (g *AtlasGlossary) Creator(name string, icon string) {
	entity := assets.Glossary{
		TypeName: "AtlasGlossary",
		Attributes: assets.GlossaryAttributes{
			Name:          name,
			QualifiedName: name,
			AssetIcon:     icon,
		},
	}
	if icon != "" {
		entity.Attributes.AssetIcon = icon
	}

	g.Entities = append(g.Entities, entity)
}

// Updater is used to modify a glossary asset in memory.
func (g *AtlasGlossary) Updater(name string, qualifiedName string, glossary_guid string) error {
	if name == "" || qualifiedName == "" || glossary_guid == "" {
		return errors.New("name, qualified_name, and glossary_guid are required fields")
	}

	entity := assets.Glossary{
		TypeName: "AtlasGlossary",
		Attributes: assets.GlossaryAttributes{
			Name:          name,
			QualifiedName: qualifiedName,
		},
		Guid: glossary_guid,
	}
	g.Entities = append(g.Entities, entity)
	return nil
}

// MarshalJSON filters out entities to only include those with non-empty attributes.
func (g *AtlasGlossary) MarshalJSON() ([]byte, error) {
	// Filter out entities to only include those with non-empty attributes
	filteredEntities := make([]assets.Glossary, 0)
	for _, entity := range g.Entities {
		if entity.Attributes.Name != "" || entity.Attributes.QualifiedName != "" || entity.Attributes.AssetIcon != "" {
			filteredEntities = append(filteredEntities, entity)
		}
	}

	type Alias AtlasGlossary

	customJSON := &struct {
		Entities []assets.Glossary `json:"entities"`
	}{
		Entities: filteredEntities,
	}

	return json.MarshalIndent(customJSON, "", "  ")
}
