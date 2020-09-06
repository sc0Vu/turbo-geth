package migrations

import (
	"context"
	"fmt"
	"github.com/ledgerwatch/lmdb-go/lmdb"
	"testing"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/stretchr/testify/require"
)

func TestDupSortHashState(t *testing.T) {
	require, db := require.New(t), ethdb.NewMemDatabase()

	err := db.KV().Update(context.Background(), func(tx ethdb.Tx) error {
		return tx.(ethdb.BucketMigrator).CreateBucket(dbutils.CurrentStateBucketOld1)
	})
	require.NoError(err)

	accKey := string(common.FromHex(fmt.Sprintf("%064x", 0)))
	inc := string(common.FromHex("0000000000000001"))
	storageKey := accKey + inc + accKey

	err = db.Put(dbutils.CurrentStateBucketOld1, []byte(accKey), []byte{1})
	require.NoError(err)
	err = db.Put(dbutils.CurrentStateBucketOld1, []byte(storageKey), []byte{2})
	require.NoError(err)

	migrator := NewMigrator()
	migrator.Migrations = []Migration{dupSortHashState}
	err = migrator.Apply(db, "")
	require.NoError(err)

	// test high-level data access didn't change
	i := 0
	err = db.Walk(dbutils.CurrentStateBucket, nil, 0, func(k, v []byte) (bool, error) {
		i++
		return true, nil
	})
	require.NoError(err)
	require.Equal(2, i)

	v, err := db.Get(dbutils.CurrentStateBucket, []byte(accKey))
	require.NoError(err)
	require.Equal([]byte{1}, v)

	v, err = db.Get(dbutils.CurrentStateBucket, []byte(storageKey))
	require.NoError(err)
	require.Equal([]byte{2}, v)

	// test low-level data layout
	rawKV := db.KV().(*ethdb.LmdbKV)
	env := rawKV.Env()
	allDBI := rawKV.AllDBI()

	tx, err := env.BeginTxn(nil, lmdb.Readonly)
	require.NoError(err)
	c, err := tx.OpenCursor(allDBI[dbutils.CurrentStateBucket])
	require.NoError(err)

	k, v, err := c.Get([]byte(accKey), nil, lmdb.Set)
	require.NoError(err)
	require.Equal([]byte(accKey), k)
	require.Equal([]byte{1}, v)

	keyLen := common.HashLength + common.IncarnationLength
	k, v, err = c.Get([]byte(storageKey)[:keyLen], []byte(storageKey)[keyLen:], lmdb.GetBothRange)
	require.NoError(err)
	require.Equal([]byte(storageKey)[:keyLen], k)
	require.Equal([]byte(storageKey)[keyLen:], v[:common.HashLength])
	require.Equal([]byte{2}, v[common.HashLength:])
}

func TestDupSortPlainState(t *testing.T) {
	require, db := require.New(t), ethdb.NewMemDatabase()

	err := db.KV().Update(context.Background(), func(tx ethdb.Tx) error {
		return tx.(ethdb.BucketMigrator).CreateBucket(dbutils.PlainStateBucketOld1)
	})
	require.NoError(err)

	accKey := string(common.FromHex(fmt.Sprintf("%040x", 0)))
	inc := string(common.FromHex("0000000000000001"))
	storageKey := accKey + inc + string(common.FromHex(fmt.Sprintf("%064x", 0)))

	err = db.Put(dbutils.PlainStateBucketOld1, []byte(accKey), []byte{1})
	require.NoError(err)
	err = db.Put(dbutils.PlainStateBucketOld1, []byte(storageKey), []byte{2})
	require.NoError(err)

	migrator := NewMigrator()
	migrator.Migrations = []Migration{dupSortPlainState}
	err = migrator.Apply(db, "")
	require.NoError(err)

	// test high-level data access didn't change
	i := 0
	err = db.Walk(dbutils.PlainStateBucket, nil, 0, func(k, v []byte) (bool, error) {
		i++
		return true, nil
	})
	require.NoError(err)
	require.Equal(2, i)

	v, err := db.Get(dbutils.PlainStateBucket, []byte(accKey))
	require.NoError(err)
	require.Equal([]byte{1}, v)

	v, err = db.Get(dbutils.PlainStateBucket, []byte(storageKey))
	require.NoError(err)
	require.Equal([]byte{2}, v)

	// test low-level data layout
	rawKV := db.KV().(*ethdb.LmdbKV)
	env := rawKV.Env()
	allDBI := rawKV.AllDBI()

	tx, err := env.BeginTxn(nil, lmdb.Readonly)
	require.NoError(err)
	c, err := tx.OpenCursor(allDBI[dbutils.PlainStateBucket])
	require.NoError(err)

	k, v, err := c.Get([]byte(accKey), nil, lmdb.Set)
	require.NoError(err)
	require.Equal([]byte(accKey), k)
	require.Equal([]byte{1}, v)

	keyLen := common.AddressLength + common.IncarnationLength
	k, v, err = c.Get([]byte(storageKey)[:keyLen], []byte(storageKey)[keyLen:], lmdb.GetBothRange)
	require.NoError(err)
	require.Equal([]byte(storageKey)[:keyLen], k)
	require.Equal([]byte(storageKey)[keyLen:], v[:common.HashLength])
	require.Equal([]byte{2}, v[common.HashLength:])
}

func TestDupSortIH(t *testing.T) {
	require, db := require.New(t), ethdb.NewMemDatabase()

	hash32Bytes := common.FromHex("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	err := db.KV().Update(context.Background(), func(tx ethdb.Tx) error {
		return tx.(ethdb.BucketMigrator).CreateBucket(dbutils.IntermediateTrieHashBucketOld1)
	})
	require.NoError(err)

	accKey := string(common.FromHex(fmt.Sprintf("%064x", 0)))
	inc := string(common.FromHex("0000000000000001"))
	storagePrefix := string(common.FromHex("77"))
	storagePrefix2 := string(common.FromHex("88"))

	err = db.Put(dbutils.IntermediateTrieHashBucketOld1, []byte{1}, hash32Bytes)
	require.NoError(err)
	err = db.Put(dbutils.IntermediateTrieHashBucketOld1, []byte(accKey+inc), hash32Bytes)
	require.NoError(err)
	err = db.Put(dbutils.IntermediateTrieHashBucketOld1, []byte(accKey+inc+storagePrefix), hash32Bytes)
	require.NoError(err)
	err = db.Put(dbutils.IntermediateTrieHashBucketOld1, []byte(accKey+inc+storagePrefix2), hash32Bytes)
	require.NoError(err)

	migrator := NewMigrator()
	migrator.Migrations = []Migration{dupSortIH}
	err = migrator.Apply(db, "")
	require.NoError(err)

	// test high-level data access didn't change
	i := 0
	err = db.Walk(dbutils.IntermediateTrieHashBucket, nil, 0, func(k, v []byte) (bool, error) {
		i++
		return true, nil
	})
	require.NoError(err)
	require.Equal(3, i)

	v, err := db.Get(dbutils.IntermediateTrieHashBucket, []byte(accKey))
	require.NoError(err)
	require.Equal([]byte{1}, v)

	//v, err = db.Get(dbutils.IntermediateTrieHashBucket, []byte(storageKey))
	//require.NoError(err)
	//require.Equal([]byte{2}, v)

	// test low-level data layout
	kv := db.KV()

	tx, err := kv.Begin(context.Background(), nil, false)
	require.NoError(err)
	defer tx.Rollback()

	c := tx.CursorDupSort(dbutils.IntermediateTrieHashBucket)
	v, err = c.SeekExact([]byte(accKey))
	require.NoError(err)
	require.Equal([]byte{1}, v)

	k, v, err := c.SeekBothRange([]byte(accKey+inc), []byte(accKey+inc))
	require.NoError(err)
	require.Equal(v[:1], []byte(storagePrefix))
	require.Equal(k, []byte(accKey+inc))

	k, v, err = c.Next()
	require.NoError(err)
	require.Equal(v[:1], []byte(storagePrefix2))
	require.Equal(k, []byte(accKey+inc))
}
