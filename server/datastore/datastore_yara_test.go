package datastore

import (
	"testing"

	"github.com/kolide/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testYARAStore(t *testing.T, ds kolide.Datastore) {
	ysg := &kolide.YARASignatureGroup{
		SignatureName: "sig1",
		Paths: []string{
			"path1",
			"path2",
		},
	}
	ysg, err := ds.NewYARASignatureGroup(ysg)
	require.Nil(t, err)
	require.True(t, ysg.ID > 0)
	fp := &kolide.FIMSection{
		SectionName: "fp1",
		Paths: []string{
			"path1",
			"path2",
			"path3",
		},
	}
	fp, err = ds.NewFIMSection(fp)
	require.Nil(t, err)
	assert.True(t, fp.ID > 0)

	err = ds.NewYARAFilePath("fp1", "sig1")
	require.Nil(t, err)
	yaraSection, err := ds.YARASection()
	require.Nil(t, err)
	require.Len(t, yaraSection.FilePaths, 1)
	assert.Len(t, yaraSection.FilePaths["fp1"], 1)
	require.Len(t, yaraSection.Signatures, 1)
	assert.Len(t, yaraSection.Signatures["sig1"], 2)
	ysg = &kolide.YARASignatureGroup{
		SignatureName: "sig2",
		Paths: []string{
			"path3",
		},
	}
	ysg, err = ds.NewYARASignatureGroup(ysg)
	require.Nil(t, err)
	yaraSection, err = ds.YARASection()
	require.Nil(t, err)
	assert.Len(t, yaraSection.Signatures["sig2"], 1)
}

func testYARATransactions(t *testing.T, ds kolide.Datastore) {
	if ds.Name() == "inmem" {
		t.Skip("not implemented for inmem")
	}

	ysg := &kolide.YARASignatureGroup{
		SignatureName: "sig1",
		Paths: []string{
			"path1",
			"path2",
		},
	}
	ysg, err := ds.NewYARASignatureGroup(ysg)
	require.Nil(t, err)
	require.True(t, ysg.ID > 0)
	fp := &kolide.FIMSection{
		SectionName: "fp1",
		Paths: []string{
			"path1",
			"path2",
			"path3",
		},
	}
	fp, err = ds.NewFIMSection(fp)
	require.Nil(t, err)
	assert.True(t, fp.ID > 0)
	tx, err := ds.Begin()
	require.Nil(t, err)

	err = ds.NewYARAFilePath("fp1", "sig1", kolide.HasTransaction(tx))
	require.Nil(t, err)
	err = tx.Rollback()
	require.Nil(t, err)
	yaraSection, err := ds.YARASection()
	require.Nil(t, err)
	require.NotNil(t, yaraSection)
	// there shouldn't be any file paths because we rolled back the transaciton
	require.Len(t, yaraSection.FilePaths, 0)

	// try it again
	tx, err = ds.Begin()
	require.Nil(t, err)
	err = ds.NewYARAFilePath("fp1", "sig1", kolide.HasTransaction(tx))
	require.Nil(t, err)
	err = tx.Commit()
	require.Nil(t, err)
	yaraSection, err = ds.YARASection()
	require.Nil(t, err)
	require.NotNil(t, yaraSection)
	// file path should exist because we committed
	require.Len(t, yaraSection.FilePaths, 1)

}
