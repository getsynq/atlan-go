package assets

import (
	"fmt"
	"testing"

	"github.com/getsynq/atlan-go/atlan"
	"github.com/stretchr/testify/assert"
)

var GlossaryName = atlan.MakeUnique("GLS")

func TestIntegrationGlossary(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := NewContext()

	glossaryGUID, glossaryQualifiedName := testCreateGlossary(t, client)
	fmt.Printf("glossaryQn: %v\n", glossaryQualifiedName)
	testUpdateGlossary(t, client, glossaryGUID)
	testRetrieveGlossary(t, client, glossaryGUID)
	testRetrieveGlossarybyQualifiedName(t, client, glossaryQualifiedName)
	testDeleteGlossary(t, client, glossaryGUID)
}

func testCreateGlossary(t *testing.T, client *AtlanClient) (string, string) {
	g := &AtlasGlossary{}
	// Create Glossary
	g.Creator(GlossaryName, atlan.AtlanIconAirplaneInFlight)
	response, err := Save(client, g)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, response, "fetched glossary should not be nil")
	assert.Len(t, response.MutatedEntities.CREATE, 1, "number of glossaries created should be 1")
	assert.Empty(t, response.MutatedEntities.UPDATE, "number of glossaries updated should be 0")
	assert.Empty(t, response.MutatedEntities.DELETE, "number of glossaries deleted should be 0")
	assetone := response.MutatedEntities.CREATE[0]
	assert.NotNil(t, assetone, "glossary should not be nil")
	assert.Equal(t, GlossaryName, *assetone.Attributes.Name, "glossary name should match")
	assert.Equal(t, *g.TypeName, assetone.TypeName, "glossary type should match")

	return assetone.Guid, *assetone.Attributes.QualifiedName
}

func testUpdateGlossary(t *testing.T, client *AtlanClient, glossaryGUID string) {
	g := &AtlasGlossary{}
	glossaryQualifiedName := GlossaryName + "-qual"
	DisplayName := "gsdk-test-update"
	g.Updater(GlossaryName, glossaryQualifiedName, glossaryGUID)
	g.DisplayName = &DisplayName
	updateresponse, err := Save(client, g)
	if err != nil {
		fmt.Println("Error:", err)
	}
	assert.NotNil(t, updateresponse, "fetched glossary should not be nil")
	assert.Len(t, updateresponse.MutatedEntities.UPDATE, 1, "number of glossaries updated should be 1")
	assert.Equal(t, *g.DisplayName, *updateresponse.MutatedEntities.UPDATE[0].Attributes.DisplayText, "glossary display name should match")
}

func testRetrieveGlossary(t *testing.T, client *AtlanClient, glossaryGUID string) {
	glossary, err := GetByGuid[*AtlasGlossary](client, glossaryGUID)
	if err != nil {
		fmt.Println("Error:", err)
	}
	assert.NotNil(t, glossary, "fetched glossary should be nil")
	assert.Equal(t, glossaryGUID, *glossary.Guid, "glossary guid should match")
}

func testRetrieveGlossarybyQualifiedName(t *testing.T, client *AtlanClient, glossaryQualifiedName string) {
	glossary, err := GetByQualifiedName[*AtlasGlossary](client, glossaryQualifiedName)
	if err != nil {
		fmt.Println("Error:", err)
	}
	assert.NotNil(t, glossary, "fetched glossary should not be nil")
	assert.Equal(t, glossaryQualifiedName, *glossary.QualifiedName, "glossary qualified name should match")
}

func testDeleteGlossary(t *testing.T, client *AtlanClient, glossaryGUID string) {
	deleteresponse, err := PurgeByGuid(client, []string{glossaryGUID})
	if err != nil {
		fmt.Println("Error:", err)
	}
	assert.NotNil(t, deleteresponse, "fetched glossary should not be nil")
	assert.Len(t, deleteresponse.MutatedEntities.DELETE, 1, "number of glossaries deleted should be 1")
	assert.Equal(t, glossaryGUID, deleteresponse.MutatedEntities.DELETE[0].Guid, "glossary guid should match")
}
