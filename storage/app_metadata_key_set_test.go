package storage_test

import (
	. "AppMetadataAPIServerGo/storage"
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestAppMetadataKeySet(t *testing.T) {
	keyset:=AppMetadataKeySet{}
	key1:= AppMetadataKey{"app", "version"}
	key2:= AppMetadataKey{"app", "version2"}
	key3:= AppMetadataKey{"app2", "version2"}

	keyset.Delete(key1) //no-op when delete on non-exists

	assert.Equal(t, keyset.GetAllAppMetadataKeys(), []AppMetadataKey{})

	assert.Equal(t, keyset.Exists(key1), false)
	keyset.Add(key1)
	keyset.Add(key1)
	keyset.Add(key1)
	keyset.Add(key1)
	assert.Equal(t, keyset.Exists(key1), true)
	assert.Equal(t, keyset.Size(), 1)
	assert.Equal(t, keyset.GetAllAppMetadataKeys(), []AppMetadataKey{key1})


	keyset.Add(key2)
	keyset.Add(key3)
	assert.Equal(t, keyset.Exists(key2), true)
	assert.Equal(t, keyset.Exists(key3), true)

	actualKeys:=keyset.GetAllAppMetadataKeys()
	assert.ElementsMatch(t, actualKeys, []AppMetadataKey{key1, key2, key3})


	keyset.Delete(key1)
	assert.Equal(t, keyset.Exists(key1), false)
	assert.Equal(t, keyset.Size(), 2)
	assert.ElementsMatch(t, keyset.GetAllAppMetadataKeys(), []AppMetadataKey{key2, key3})


	keyset.Delete(key2)
	keyset.Delete(key3)
	assert.Equal(t, keyset.Size(), 0)
	assert.ElementsMatch(t, keyset.GetAllAppMetadataKeys(), []AppMetadataKey{})

}

func sorted(keys []AppMetadataKey) []AppMetadataKey{
	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i].Title <= keys[j].Title && keys[i].Version <= keys[j].Version
	})
	return keys
}