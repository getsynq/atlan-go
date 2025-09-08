package assets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/getsynq/atlan-go/atlan"
	"github.com/stretchr/testify/assert"
)

var PurposeName = atlan.MakeUnique("Purpose")

func TestIntegrationPurpose(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	client := NewContext()
	// ctx.EnableLogging("debug")

	purposeID, purposeQualifiedName := testCreatePurpose(t, client)
	testRetrievePurpose(t, client, purposeID)
	testPurposeCreateMetadataPolicy(t, client, purposeID)
	testPurposeCreateDataPolicy(t, client, purposeID)
	testFindPurposesByName(t, client)
	testUpdatePurpose(t, client, purposeQualifiedName)
	testDeletePurpose(t, client, purposeID)
}

func testCreatePurpose(t *testing.T, client *AtlanClient) (string, string) {
	p := &Purpose{
		client: client,
	}
	// Create Purpose
	atlanTags := []string{"Issue", "Confidential"}
	err := p.Creator(PurposeName, atlanTags)
	require.NoError(t, err, "creator should not return an error")

	response, err := Save(client)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, response, "fetched purpose should not be nil")
	assert.Len(t, response.MutatedEntities.CREATE, 1, "number of purposes created should be 1")
	assert.Empty(t, response.MutatedEntities.UPDATE, "number of purposes updated should be 0")
	assert.Empty(t, response.MutatedEntities.DELETE, "number of purposes deleted should be 0")
	CreatedPurpose := response.MutatedEntities.CREATE[0]
	assert.NotNil(t, CreatedPurpose, "purpose should not be nil")
	assert.Equal(t, PurposeName, *CreatedPurpose.Attributes.Name, "purpose name should match")
	assert.Equal(t, *p.TypeName, CreatedPurpose.TypeName, "purpose type should match")

	return CreatedPurpose.Guid, *CreatedPurpose.Attributes.QualifiedName
}

func testRetrievePurpose(t *testing.T, client *AtlanClient, purposeID string) {
	purpose, err := GetByGuid[*Purpose](client, purposeID)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, purpose, "fetched purpose should not be nil")
	assert.Equal(t, PurposeName, *purpose.Name, "purpose name should match")
}

func testFindPurposesByName(t *testing.T, client *AtlanClient) {
	time.Sleep(3 * time.Second)
	purposes, err := FindPurposesByName(PurposeName)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, purposes.Entities, "fetched purposes should not be nil")
	assert.Equal(t, int64(1), purposes.ApproximateCount, "number of purposes fetched should be 1")
	assert.Equal(t, PurposeName, *purposes.Entities[0].Name, "purpose name should match")
}

func testPurposeCreateMetadataPolicy(t *testing.T, client *AtlanClient, purposeID string) {
	p := &Purpose{}
	policy, err := p.CreateMetadataPolicy(
		PurposeName,
		purposeID,
		atlan.AuthPolicyTypeAllow,
		[]atlan.PurposeMetadataAction{
			atlan.PurposeMetadataActionRead,
		},
		nil,
		nil,
		true,
	)
	require.NoError(t, err, "error should be nil while creating metadata policy")
	response, err := Save(client, policy)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, response, "fetched policy should not be nil")
	assert.Len(t, response.MutatedEntities.CREATE, 1, "number of policies added should be 1")
	CreatedPolicy := response.MutatedEntities.CREATE[0]
	assert.NotNil(t, CreatedPolicy, "policy should not be nil")
}

func testPurposeCreateDataPolicy(t *testing.T, client *AtlanClient, purposeID string) {
	p := &Purpose{}
	policy, err := p.CreateDataPolicy(
		PurposeName,
		purposeID,
		atlan.AuthPolicyTypeAllow,
		nil,
		nil,
		true,
	)
	require.NoError(t, err, "error should be nil while creating data policy")
	response, err := Save(client, policy)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, response, "fetched policy should not be nil")
	assert.Len(t, response.MutatedEntities.CREATE, 1, "number of policies added should be 1")
	CreatedPolicy := response.MutatedEntities.CREATE[0]
	assert.NotNil(t, CreatedPolicy, "policy should not be nil")
}

func testUpdatePurpose(t *testing.T, client *AtlanClient, purposeQualifiedName string) {
	p := &Purpose{}
	NewName := atlan.MakeUnique("test-update-purpose")
	Description := atlan.MakeUnique("test-update-purpose-description")
	err := p.Updater(purposeQualifiedName, PurposeName, true)
	require.NoError(t, err, "updater should not return an error")

	p.Name = &NewName
	p.Description = &Description
	UpdaterResponse, err := Save(client, p)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	assert.NotNil(t, UpdaterResponse, "fetched purpose should not be nil")
	assert.Len(t, UpdaterResponse.MutatedEntities.UPDATE, 1, "number of purposes updated should be 1")
	assert.Equal(t, *p.Name, *UpdaterResponse.MutatedEntities.UPDATE[0].Attributes.Name, "purpose name should match")
	assert.Equal(t, *p.Description, *UpdaterResponse.MutatedEntities.UPDATE[0].Attributes.Description, "purpose description should match")
}

func testDeletePurpose(t *testing.T, client *AtlanClient, purposeID string) {
	DeleteResponse, err := PurgeByGuid(client, []string{purposeID})
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	// for _, deleted := range DeleteResponse.MutatedEntities.DELETE {
	//	t.Logf("Deleted: %v", deleted)
	//}
	assert.NotNil(t, DeleteResponse, "fetched purpose should not be nil")
	assert.Len(t, DeleteResponse.MutatedEntities.DELETE, 3, "number of purposes deleted should be 3") // 3 because of the metadata and data policies

	// Collect GUIDs from the server response
	serverGuids := make([]string, len(DeleteResponse.MutatedEntities.DELETE))
	for i, deleted := range DeleteResponse.MutatedEntities.DELETE {
		serverGuids[i] = deleted.Guid
	}

	// Ensure the expected purposeID is in the list of server-provided GUIDs
	assert.Contains(t, serverGuids, purposeID, "purpose guid should match one of the server-provided GUIDs")
}
