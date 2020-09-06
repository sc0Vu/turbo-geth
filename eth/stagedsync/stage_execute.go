package stagedsync

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/core/types/accounts"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"github.com/ledgerwatch/turbo-geth/eth/stagedsync/stages"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/params"
	"github.com/ledgerwatch/turbo-geth/rlp"
)

const (
	logInterval = 30 * time.Second
)

type HasChangeSetWriter interface {
	ChangeSetWriter() *state.ChangeSetWriter
}

type ChangeSetHook func(blockNum uint64, wr *state.ChangeSetWriter)

func SpawnExecuteBlocksStage(s *StageState, stateDB ethdb.Database, chainConfig *params.ChainConfig, chainContext core.ChainContext, vmConfig *vm.Config, toBlock uint64, quit <-chan struct{}, writeReceipts bool, hdd bool, changeSetHook ChangeSetHook) error {
	prevStageProgress, _, errStart := stages.GetStageProgress(stateDB, stages.Senders)
	if errStart != nil {
		return errStart
	}
	var to = prevStageProgress
	if toBlock > 0 {
		to = min(prevStageProgress, toBlock)
	}
	if to <= s.BlockNumber {
		s.Done()
		return nil
	}
	log.Info("Blocks execution", "from", s.BlockNumber, "to", to)

	if prof {
		f, err := os.Create(fmt.Sprintf("cpu-%d.prof", s.BlockNumber))
		if err != nil {
			log.Error("could not create CPU profile", "error", err)
			return err
		}
		if err = pprof.StartCPUProfile(f); err != nil {
			log.Error("could not start CPU profile", "error", err)
			return err
		}
	}

	var tx ethdb.DbWithPendingMutations
	var useExternalTx bool
	if hasTx, ok := stateDB.(ethdb.HasTx); ok && hasTx.Tx() != nil {
		tx = stateDB.(ethdb.DbWithPendingMutations)
		useExternalTx = true
	} else {
		var err error
		tx, err = stateDB.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
	}

	batch := tx.NewBatch()
	defer batch.Rollback()

	engine := chainContext.Engine()

	stageProgress := s.BlockNumber
	logEvery := time.NewTicker(logInterval)
	defer logEvery.Stop()
	logBlock := stageProgress
	// Warmup only works for HDD sync, and for long ranges
	var warmup = hdd && (to-s.BlockNumber) > 30000

	for blockNum := stageProgress + 1; blockNum <= to; blockNum++ {
		if err := common.Stopped(quit); err != nil {
			return err
		}

		stageProgress = blockNum

		blockHash := rawdb.ReadCanonicalHash(tx, blockNum)
		block := rawdb.ReadBlock(tx, blockHash, blockNum)
		if block == nil {
			break
		}
		senders := rawdb.ReadSenders(tx, blockHash, blockNum)
		block.Body().SendersToTxs(senders)

		if warmup {
			log.Info("Running a warmup...")
			count := 0
			if err := stateDB.Walk(dbutils.PlainStateBucket, nil, 0, func(_, _ []byte) (bool, error) {
				if err := common.Stopped(quit); err != nil {return false, nil
				}
				count++
				if count%10000000 == 0 {
					log.Info("Warmed up", "keys", count)
				}
				return true, nil
			}); err != nil {
				return err
			}
			warmup = false
			log.Info("Warm up done.")
		}

		var stateReader state.StateReader
		var stateWriter state.WriterWithChangeSets

		stateReader = state.NewPlainStateReader(batch)
		stateWriter = state.NewPlainStateWriter(batch, tx, blockNum)

		// where the magic happens
		receipts, err := core.ExecuteBlockEphemerally(chainConfig, vmConfig, chainContext, engine, block, stateReader, stateWriter)
		if err != nil {
			return err
		}

		if writeReceipts {
			// Convert the receipts into their storage form and serialize them
			storageReceipts := make([]*types.ReceiptForStorage, len(receipts))
			for i, receipt := range receipts {
				storageReceipts[i] = (*types.ReceiptForStorage)(receipt)
			}
			var bytes []byte
			if bytes, err = rlp.EncodeToBytes(storageReceipts); err != nil {
				return fmt.Errorf("encode block receipts for block %d: %v", block.NumberU64(), err)
			}
			// Store the flattened receipt slice
			if err = tx.Append(dbutils.BlockReceiptsPrefix, dbutils.BlockReceiptsKey(block.NumberU64(), block.Hash()), bytes); err != nil {
				return fmt.Errorf("writing receipts for block %d: %v", block.NumberU64(), err)
			}
		}

		if batch.BatchSize() >= batch.IdealBatchSize() {
			if err = s.Update(batch, blockNum); err != nil {
				return err
			}
			if err = batch.CommitAndBegin(); err != nil {
				return err
			}
			if !useExternalTx {
				if err = tx.CommitAndBegin(); err != nil {
					return err
				}
			}
			warmup = hdd && (to-blockNum) > 30000
		}

		if prof {
			if blockNum-s.BlockNumber == 100000 {
				// Flush the CPU profiler
				pprof.StopCPUProfile()
			}
		}

		if changeSetHook != nil {
			if hasChangeSet, ok := stateWriter.(HasChangeSetWriter); ok {
				changeSetHook(blockNum, hasChangeSet.ChangeSetWriter())
			}
		}

		select {
		default:
		case <-logEvery.C:
			logBlock = logProgress(logBlock, blockNum, batch)
		}
	}

	if err := s.Update(batch, stageProgress); err != nil {
		return err
	}
	if _, err := batch.Commit(); err != nil {
		return fmt.Errorf("sync Execute: failed to write batch commit: %v", err)
	}
	if !useExternalTx {
		if _, err := tx.Commit(); err != nil {
			return err
		}
	}
	log.Info("Completed on", "block", stageProgress)
	s.Done()
	return nil
}

func logProgress(prev, now uint64, batch ethdb.DbWithPendingMutations) uint64 {
	speed := float64(now-prev) / float64(logInterval/time.Second)
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Info("Executed blocks:",
		"currentBlock", now,
		"blk/second", speed,
		"batch", common.StorageSize(batch.BatchSize()),
		"alloc", common.StorageSize(m.Alloc),
		"sys", common.StorageSize(m.Sys),
		"numGC", int(m.NumGC))

	return now
}

func UnwindExecutionStage(u *UnwindState, s *StageState, stateDB ethdb.Database, writeReceipts bool) error {
	if u.UnwindPoint >= s.BlockNumber {
		s.Done()
		return nil
	}

	log.Info("Unwind Execution stage", "from", s.BlockNumber, "to", u.UnwindPoint)
	batch := stateDB.NewBatch()
	defer batch.Rollback()

	rewindFunc := ethdb.RewindDataPlain
	stateBucket := dbutils.PlainStateBucket
	storageKeyLength := common.AddressLength + common.IncarnationLength + common.HashLength
	deleteAccountFunc := deleteAccountPlain
	writeAccountFunc := writeAccountPlain
	recoverCodeHashFunc := recoverCodeHashPlain

	accountMap, storageMap, err := rewindFunc(stateDB, s.BlockNumber, u.UnwindPoint)
	if err != nil {
		return fmt.Errorf("unwind Execution: getting rewind data: %v", err)
	}

	for key, value := range accountMap {
		if len(value) > 0 {
			var acc accounts.Account
			if err = acc.DecodeForStorage(value); err != nil {
				return err
			}

			// Fetch the code hash
			recoverCodeHashFunc(&acc, stateDB, key)
			if err = writeAccountFunc(batch, key, acc); err != nil {
				return err
			}
		} else {
			if err = deleteAccountFunc(batch, key); err != nil {
				return err
			}
		}
	}
	for key, value := range storageMap {
		if len(value) > 0 {
			if err = batch.Put(stateBucket, []byte(key)[:storageKeyLength], value); err != nil {
				return err
			}
		} else {
			if err = batch.Delete(stateBucket, []byte(key)[:storageKeyLength]); err != nil {
				return err
			}
		}
	}

	if err = stateDB.Walk(dbutils.PlainAccountChangeSetBucket, dbutils.EncodeTimestamp(u.UnwindPoint+1), 0, func(k, _ []byte) (bool, error) {
		if err1 := batch.Delete(dbutils.PlainAccountChangeSetBucket, common.CopyBytes(k)); err1 != nil {
			return false, fmt.Errorf("unwind Execution: delete account changesets: %v", err1)
		}
		return true, nil
	}); err != nil {
		return fmt.Errorf("unwind Execution: walking account changesets: %v", err)
	}
	if err = stateDB.Walk(dbutils.PlainStorageChangeSetBucket, dbutils.EncodeTimestamp(u.UnwindPoint+1), 0, func(k, _ []byte) (bool, error) {
		if err1 := batch.Delete(dbutils.PlainStorageChangeSetBucket, common.CopyBytes(k)); err1 != nil {
			return false, fmt.Errorf("unwind Execution: delete storage changesets: %v", err1)
		}
		return true, nil
	}); err != nil {
		return fmt.Errorf("unwind Execution: walking storage changesets: %v", err)
	}
	if writeReceipts {
		if err = stateDB.Walk(dbutils.BlockReceiptsPrefix, dbutils.EncodeBlockNumber(u.UnwindPoint+1), 0, func(k, _ []byte) (bool, error) {
			if err1 := batch.Delete(dbutils.BlockReceiptsPrefix, common.CopyBytes(k)); err1 != nil {
				return false, fmt.Errorf("unwind Execution: delete receipts: %v", err1)
			}
			return true, nil
		}); err != nil {
			return fmt.Errorf("unwind Execution: walking receipts: %v", err)
		}
	}

	if err = u.Done(batch); err != nil {
		return fmt.Errorf("unwind Execution: reset: %v", err)
	}

	_, err = batch.Commit()
	if err != nil {
		return fmt.Errorf("unwind Execute: failed to write db commit: %v", err)
	}
	return nil
}

func writeAccountHashed(db ethdb.Database, key string, acc accounts.Account) error {
	var addrHash common.Hash
	copy(addrHash[:], []byte(key))
	if err := cleanupContractCodeBucket(
		db,
		dbutils.ContractCodeBucket,
		acc,
		func(db ethdb.Getter, out *accounts.Account) (bool, error) {
			return rawdb.ReadAccount(db, addrHash, out)
		},
		func(inc uint64) []byte { return dbutils.GenerateStoragePrefix(addrHash[:], inc) },
	); err != nil {
		return err
	}
	return rawdb.WriteAccount(db, addrHash, acc)
}

func writeAccountPlain(db ethdb.Database, key string, acc accounts.Account) error {
	var address common.Address
	copy(address[:], []byte(key))
	if err := cleanupContractCodeBucket(
		db,
		dbutils.PlainContractCodeBucket,
		acc,
		func(db ethdb.Getter, out *accounts.Account) (bool, error) {
			return rawdb.PlainReadAccount(db, address, out)
		},
		func(inc uint64) []byte { return dbutils.PlainGenerateStoragePrefix(address[:], inc) },
	); err != nil {
		return fmt.Errorf("writeAccountPlain for %x: %w", address, err)
	}

	return rawdb.PlainWriteAccount(db, address, acc)
}

func recoverCodeHashHashed(acc *accounts.Account, db ethdb.Getter, key string) {
	var addrHash common.Hash
	copy(addrHash[:], []byte(key))
	if acc.Incarnation > 0 && acc.IsEmptyCodeHash() {
		if codeHash, err2 := db.Get(dbutils.ContractCodeBucket, dbutils.GenerateStoragePrefix(addrHash[:], acc.Incarnation)); err2 == nil {
			copy(acc.CodeHash[:], codeHash)
		}
	}
}

func cleanupContractCodeBucket(
	db ethdb.Database,
	bucket string,
	acc accounts.Account,
	readAccountFunc func(ethdb.Getter, *accounts.Account) (bool, error),
	getKeyForIncarnationFunc func(uint64) []byte,
) error {
	var original accounts.Account
	got, err := readAccountFunc(db, &original)
	if err != nil && !errors.Is(err, ethdb.ErrKeyNotFound) {
		return fmt.Errorf("cleanupContractCodeBucket: %w", err)
	}
	if got {
		// clean up all the code incarnations original incarnation and the new one
		for incarnation := original.Incarnation; incarnation > acc.Incarnation && incarnation > 0; incarnation-- {
			err = db.Delete(bucket, getKeyForIncarnationFunc(incarnation))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func recoverCodeHashPlain(acc *accounts.Account, db ethdb.Getter, key string) {
	var address common.Address
	copy(address[:], []byte(key))
	if acc.Incarnation > 0 && acc.IsEmptyCodeHash() {
		if codeHash, err2 := db.Get(dbutils.PlainContractCodeBucket, dbutils.PlainGenerateStoragePrefix(address[:], acc.Incarnation)); err2 == nil {
			copy(acc.CodeHash[:], codeHash)
		}
	}
}

func deleteAccountHashed(db rawdb.DatabaseDeleter, key string) error {
	var addrHash common.Hash
	copy(addrHash[:], []byte(key))
	return rawdb.DeleteAccount(db, addrHash)
}

func deleteAccountPlain(db rawdb.DatabaseDeleter, key string) error {
	var address common.Address
	copy(address[:], []byte(key))
	return rawdb.PlainDeleteAccount(db, address)
}

func deleteChangeSets(batch ethdb.Deleter, timestamp uint64, accountBucket, storageBucket string) error {
	changeSetKey := dbutils.EncodeTimestamp(timestamp)
	if err := batch.Delete(accountBucket, changeSetKey); err != nil {
		return err
	}
	if err := batch.Delete(storageBucket, changeSetKey); err != nil {
		return err
	}
	return nil
}

func min(a, b uint64) uint64 {
	if a <= b {
		return a
	}
	return b
}
