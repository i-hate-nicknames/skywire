package transport

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTransportLogStore(t *testing.T, logStore LogStore) {
	t.Helper()

	id1 := uuid.New()
	entry1 := &LogEntry{big.NewInt(100), big.NewInt(200)}
	id2 := uuid.New()
	entry2 := &LogEntry{big.NewInt(300), big.NewInt(400)}

	require.NoError(t, logStore.Record(id1, entry1))
	require.NoError(t, logStore.Record(id2, entry2))

	entry, err := logStore.Entry(id2)
	require.NoError(t, err)
	assert.Equal(t, int64(300), entry.ReceivedBytes.Int64())
	assert.Equal(t, int64(400), entry.SentBytes.Int64())
}

func TestInMemoryTransportLogStore(t *testing.T) {
	testTransportLogStore(t, InMemoryTransportLogStore())
}

func TestFileTransportLogStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "log_store")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	testTransportLogStore(t, FileTransportLogStore(dir))
}
