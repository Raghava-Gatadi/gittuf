// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"testing"

	"github.com/gittuf/gittuf/internal/tuf"
	"github.com/stretchr/testify/assert"
)

func TestInitializeRootMetadata(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)
	assert.Equal(t, key, rootMetadata.Keys[key.KeyID])
	assert.Equal(t, 1, rootMetadata.Roles[RootRoleName].Threshold)
	assert.Equal(t, []string{key.KeyID}, rootMetadata.Roles[RootRoleName].KeyIDs)
}

func TestAddRootKey(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	newRootKey, err := tuf.LoadKeyFromBytes(targets1KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata = AddRootKey(rootMetadata, newRootKey)

	assert.Equal(t, newRootKey, rootMetadata.Keys[newRootKey.KeyID])
	assert.Equal(t, []string{key.KeyID, newRootKey.KeyID}, rootMetadata.Roles[RootRoleName].KeyIDs)
}

func TestRemoveRootKey(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	newRootKey, err := tuf.LoadKeyFromBytes(targets1KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata = AddRootKey(rootMetadata, newRootKey)

	rootMetadata, err = DeleteRootKey(rootMetadata, newRootKey.KeyID)

	assert.Nil(t, err)
	assert.Equal(t, key, rootMetadata.Keys[key.KeyID])
	assert.Equal(t, newRootKey, rootMetadata.Keys[newRootKey.KeyID])
	assert.Equal(t, []string{key.KeyID}, rootMetadata.Roles[RootRoleName].KeyIDs)

	rootMetadata, err = DeleteRootKey(rootMetadata, key.KeyID)

	assert.ErrorIs(t, err, ErrCannotMeetThreshold)
	assert.Nil(t, rootMetadata)
}

func TestAddTargetsKey(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	targetsKey, err := tuf.LoadKeyFromBytes(targets1KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	_, err = AddTargetsKey(nil, targetsKey)
	assert.ErrorIs(t, err, ErrRootMetadataNil)

	_, err = AddTargetsKey(rootMetadata, nil)
	assert.ErrorIs(t, err, ErrTargetsKeyNil)

	rootMetadata, err = AddTargetsKey(rootMetadata, targetsKey)
	assert.Nil(t, err)
	assert.Equal(t, targetsKey, rootMetadata.Keys[targetsKey.KeyID])
	assert.Equal(t, []string{targetsKey.KeyID}, rootMetadata.Roles[TargetsRoleName].KeyIDs)
}

func TestDeleteTargetsKey(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	targetsKey1, err := tuf.LoadKeyFromBytes(targets1KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	targetsKey2, err := tuf.LoadKeyFromBytes(targets2KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata, err = AddTargetsKey(rootMetadata, targetsKey1)
	assert.Nil(t, err)
	rootMetadata, err = AddTargetsKey(rootMetadata, targetsKey2)
	assert.Nil(t, err)

	_, err = DeleteTargetsKey(nil, targetsKey1.KeyID)
	assert.ErrorIs(t, err, ErrRootMetadataNil)

	_, err = DeleteTargetsKey(rootMetadata, "")
	assert.ErrorIs(t, err, ErrKeyIDEmpty)

	rootMetadata, err = DeleteTargetsKey(rootMetadata, targetsKey1.KeyID)
	assert.Nil(t, err)
	assert.Equal(t, targetsKey1, rootMetadata.Keys[targetsKey1.KeyID])
	assert.Equal(t, targetsKey2, rootMetadata.Keys[targetsKey2.KeyID])
	targetsRole := rootMetadata.Roles[TargetsRoleName]
	assert.Contains(t, targetsRole.KeyIDs, targetsKey2.KeyID)

	rootMetadata, err = DeleteTargetsKey(rootMetadata, targetsKey2.KeyID)
	assert.ErrorIs(t, err, ErrCannotMeetThreshold)
	assert.Nil(t, rootMetadata)
}

func TestAddGitHubAppKey(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	appKey, err := tuf.LoadKeyFromBytes(targets1KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	_, err = AddGitHubAppKey(nil, appKey)
	assert.ErrorIs(t, err, ErrRootMetadataNil)

	_, err = AddGitHubAppKey(rootMetadata, nil)
	assert.ErrorIs(t, err, ErrGitHubAppKeyNil)

	rootMetadata, err = AddGitHubAppKey(rootMetadata, appKey)
	assert.Nil(t, err)
	assert.Equal(t, appKey, rootMetadata.Keys[appKey.KeyID])
	assert.Equal(t, []string{appKey.KeyID}, rootMetadata.Roles[GitHubAppRoleName].KeyIDs)
}

func TestDeleteGitHubAppKey(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	appKey, err := tuf.LoadKeyFromBytes(targets1KeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata, err = AddGitHubAppKey(rootMetadata, appKey)
	assert.Nil(t, err)

	_, err = DeleteGitHubAppKey(nil)
	assert.ErrorIs(t, err, ErrRootMetadataNil)

	rootMetadata, err = DeleteGitHubAppKey(rootMetadata)
	assert.Nil(t, err)

	assert.Empty(t, rootMetadata.Roles[GitHubAppRoleName].KeyIDs)
}

func TestEnableGitHubAppApprovals(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	_, err = EnableGitHubAppApprovals(nil)
	assert.ErrorIs(t, err, ErrRootMetadataNil)

	rootMetadata, err = EnableGitHubAppApprovals(rootMetadata)
	assert.Nil(t, err)

	assert.True(t, rootMetadata.GitHubApprovalsTrusted)
}

func TestDisableGitHubAppApprovals(t *testing.T) {
	key, err := tuf.LoadKeyFromBytes(rootKeyBytes)
	if err != nil {
		t.Fatal(err)
	}

	rootMetadata := InitializeRootMetadata(key)

	_, err = DisableGitHubAppApprovals(nil)
	assert.ErrorIs(t, err, ErrRootMetadataNil)

	rootMetadata, err = DisableGitHubAppApprovals(rootMetadata)
	assert.Nil(t, err)

	assert.False(t, rootMetadata.GitHubApprovalsTrusted)
}
