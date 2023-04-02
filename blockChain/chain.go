package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"os"
	"runtime"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"                   // DB가 존재하는지 확인하는 값
	genesisData = "this is First Transaction From CoinBase" // 첫번쨰 Genesis를 생성하기 위한 값
)

var key = []byte("lh")

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	// DB가 만들어져 잇는지 확인하는 함수
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func ContinueBlockChain(Address string) *BlockChain {

	if !DBexists() {
		fmt.Println("BlockChain Db is not Exists")
		runtime.Goexit()
	}

	var lastHash []byte

	db := getDB()

	err := db.Update(func(tx *badger.Txn) error {

		item, err := tx.Get(key)

		ErrorHandle(err)

		lastHash, err = item.ValueCopy(key)

		return err
	})

	ErrorHandle(err)
	return &BlockChain{lastHash, db}
}

func InitBlockChain(address string) *BlockChain {
	var lastHash []byte

	if DBexists() {
		fmt.Println("BlockChain Db is already Exists")
		runtime.Goexit()
	}

	db := getDB()

	err := db.Update(func(tx *badger.Txn) error {
		genesisTx := GenesisTx(address, genesisData)

		fmt.Println(genesisTx.Outputs)
		genesis := Genesis(genesisTx)

		fmt.Println(" Genesis Tx Created")

		err := tx.Set(genesis.Hash, genesis.Serialize())
		err = tx.Set(key, genesis.Hash)
		ErrorHandle(err)

		lastHash = genesis.Hash

		return err

	})
	ErrorHandle(err)

	return &BlockChain{LastHash: lastHash, Database: db}
}

func (chain *BlockChain) AddBlock(txs []*Transaction) {
	var lastHash []byte
	err := chain.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)
		ErrorHandle(err)
		lastHash, err = item.ValueCopy(key)
		return err
	})

	ErrorHandle(err)

	newBlock := CreateBlock(txs, lastHash)

	err = chain.Database.Update(func(tx *badger.Txn) error {
		err := tx.Set(newBlock.Hash, newBlock.Serialize())

		ErrorHandle(err)

		err = tx.Set(key, newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})

	ErrorHandle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(iter.CurrentHash)

		ErrorHandle(err)
		encodedBlock, err := item.ValueCopy(key)

		block = Deserialize(encodedBlock)

		return err
	})

	ErrorHandle(err)

	return block
}

func (iter *BlockChainIterator) GetByPrevHash(prevHash []byte) *Block {
	var block *Block

	err := iter.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(prevHash)

		if err != nil {
			return nil
		}

		encodedBlock, err := item.ValueCopy(key)

		block = Deserialize(encodedBlock)

		return err
	})

	if err != nil {
		return nil
	}

	return block
}

func getDB() *badger.DB {
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)

	ErrorHandle(err)

	return db
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction

	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()
	block := iter.Next()

	for {

		for _, tx := range block.Transaction {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if tx.IsGenesisTx() == false {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Out)
					}
				}
			}
		}

		block = iter.GetByPrevHash(block.PrevHash)

		if block == nil {
			break
		}

	}

	return unspentTxs
}

func (chain *BlockChain) FindUTXO(address string) []TxOutput {
	var UTXOs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
