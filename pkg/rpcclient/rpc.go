package rpcclient

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nspcc-dev/neo-go/pkg/config"
	"github.com/nspcc-dev/neo-go/pkg/config/netmode"
	"github.com/nspcc-dev/neo-go/pkg/core/block"
	"github.com/nspcc-dev/neo-go/pkg/core/fee"
	"github.com/nspcc-dev/neo-go/pkg/core/native/nativenames"
	"github.com/nspcc-dev/neo-go/pkg/core/native/nativeprices"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
	"github.com/nspcc-dev/neo-go/pkg/core/transaction"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/encoding/address"
	"github.com/nspcc-dev/neo-go/pkg/encoding/fixedn"
	"github.com/nspcc-dev/neo-go/pkg/io"
	"github.com/nspcc-dev/neo-go/pkg/neorpc"
	"github.com/nspcc-dev/neo-go/pkg/neorpc/result"
	"github.com/nspcc-dev/neo-go/pkg/network/payload"
	"github.com/nspcc-dev/neo-go/pkg/rpcclient/unwrap"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/trigger"
	"github.com/nspcc-dev/neo-go/pkg/util"
	"github.com/nspcc-dev/neo-go/pkg/vm/opcode"
	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
	"github.com/nspcc-dev/neo-go/pkg/wallet"
)

var errNetworkNotInitialized = errors.New("RPC client network is not initialized")

// CalculateNetworkFee calculates network fee for the transaction. The transaction may
// have empty witnesses for contract signers and may have only verification scripts
// filled for standard sig/multisig signers.
func (c *Client) CalculateNetworkFee(tx *transaction.Transaction) (int64, error) {
	var (
		params = []any{tx.Bytes()}
		resp   = new(result.NetworkFee)
	)
	if err := c.performRequest("calculatenetworkfee", params, resp); err != nil {
		return 0, err
	}
	return resp.Value, nil
}

// GetApplicationLog returns a contract log based on the specified txid.
func (c *Client) GetApplicationLog(hash util.Uint256, trig *trigger.Type) (*result.ApplicationLog, error) {
	var (
		params = []any{hash.StringLE()}
		resp   = new(result.ApplicationLog)
	)
	if trig != nil {
		params = append(params, trig.String())
	}
	if err := c.performRequest("getapplicationlog", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBestBlockHash returns the hash of the tallest block in the blockchain.
func (c *Client) GetBestBlockHash() (util.Uint256, error) {
	var resp = util.Uint256{}
	if err := c.performRequest("getbestblockhash", nil, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetBlockCount returns the number of blocks in the blockchain.
func (c *Client) GetBlockCount() (uint32, error) {
	var resp uint32
	if err := c.performRequest("getblockcount", nil, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetBlockByIndex returns a block by its height. In-header stateroot option
// must be initialized with Init before calling this method.
func (c *Client) GetBlockByIndex(index uint32) (*block.Block, error) {
	return c.getBlock(index)
}

// GetBlockByHash returns a block by its hash. In-header stateroot option
// must be initialized with Init before calling this method.
func (c *Client) GetBlockByHash(hash util.Uint256) (*block.Block, error) {
	return c.getBlock(hash.StringLE())
}

func (c *Client) getBlock(param any) (*block.Block, error) {
	var (
		resp []byte
		err  error
		b    *block.Block
	)
	if err = c.performRequest("getblock", []any{param}, &resp); err != nil {
		return nil, err
	}
	r := io.NewBinReaderFromBuf(resp)
	sr, err := c.StateRootInHeader()
	if err != nil {
		return nil, err
	}
	b = block.New(sr)
	b.DecodeBinary(r)
	if r.Err != nil {
		return nil, r.Err
	}
	return b, nil
}

// GetBlockByIndexVerbose returns a block wrapper with additional metadata by
// its height. In-header stateroot option must be initialized with Init before
// calling this method.
// NOTE: to get transaction.ID and transaction.Size, use t.Hash() and io.GetVarSize(t) respectively.
func (c *Client) GetBlockByIndexVerbose(index uint32) (*result.Block, error) {
	return c.getBlockVerbose(index)
}

// GetBlockByHashVerbose returns a block wrapper with additional metadata by
// its hash. In-header stateroot option must be initialized with Init before
// calling this method.
func (c *Client) GetBlockByHashVerbose(hash util.Uint256) (*result.Block, error) {
	return c.getBlockVerbose(hash.StringLE())
}

func (c *Client) getBlockVerbose(param any) (*result.Block, error) {
	var (
		params = []any{param, 1} // 1 for verbose.
		resp   = &result.Block{}
		err    error
	)
	sr, err := c.StateRootInHeader()
	if err != nil {
		return nil, err
	}
	resp.Header.StateRootEnabled = sr
	if err = c.performRequest("getblock", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBlockHash returns the hash value of the corresponding block based on the specified index.
func (c *Client) GetBlockHash(index uint32) (util.Uint256, error) {
	var (
		params = []any{index}
		resp   = util.Uint256{}
	)
	if err := c.performRequest("getblockhash", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetBlockHeader returns the corresponding block header information from a serialized hex string
// according to the specified script hash. In-header stateroot option must be
// initialized with Init before calling this method.
func (c *Client) GetBlockHeader(hash util.Uint256) (*block.Header, error) {
	var (
		params = []any{hash.StringLE()}
		resp   []byte
		h      *block.Header
	)
	if err := c.performRequest("getblockheader", params, &resp); err != nil {
		return nil, err
	}
	sr, err := c.StateRootInHeader()
	if err != nil {
		return nil, err
	}
	r := io.NewBinReaderFromBuf(resp)
	h = new(block.Header)
	h.StateRootEnabled = sr
	h.DecodeBinary(r)
	if r.Err != nil {
		return nil, r.Err
	}
	return h, nil
}

// GetBlockHeaderCount returns the number of headers in the main chain.
func (c *Client) GetBlockHeaderCount() (uint32, error) {
	var resp uint32
	if err := c.performRequest("getblockheadercount", nil, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetBlockHeaderVerbose returns the corresponding block header information from a Json format string
// according to the specified script hash. In-header stateroot option must be
// initialized with Init before calling this method.
func (c *Client) GetBlockHeaderVerbose(hash util.Uint256) (*result.Header, error) {
	var (
		params = []any{hash.StringLE(), 1}
		resp   = &result.Header{}
	)
	if err := c.performRequest("getblockheader", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBlockSysFee returns the system fees of the block based on the specified index.
// This method is only supported by NeoGo servers.
func (c *Client) GetBlockSysFee(index uint32) (fixedn.Fixed8, error) {
	var (
		params = []any{index}
		resp   fixedn.Fixed8
	)
	if err := c.performRequest("getblocksysfee", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetConnectionCount returns the current number of the connections for the node.
func (c *Client) GetConnectionCount() (int, error) {
	var resp int

	if err := c.performRequest("getconnectioncount", nil, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetCommittee returns the current public keys of NEO nodes in the committee.
func (c *Client) GetCommittee() (keys.PublicKeys, error) {
	var resp = new(keys.PublicKeys)

	if err := c.performRequest("getcommittee", nil, resp); err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetContractStateByHash queries contract information according to the contract script hash.
func (c *Client) GetContractStateByHash(hash util.Uint160) (*state.Contract, error) {
	return c.getContractState(hash.StringLE())
}

// GetContractStateByAddressOrName queries contract information using the contract
// address or name. Notice that name-based queries work only for native contracts,
// non-native ones can't be requested this way.
func (c *Client) GetContractStateByAddressOrName(addressOrName string) (*state.Contract, error) {
	return c.getContractState(addressOrName)
}

// GetContractStateByID queries contract information according to the contract ID.
// Notice that this is supported by all servers only for native contracts,
// non-native ones can be requested only from NeoGo servers.
func (c *Client) GetContractStateByID(id int32) (*state.Contract, error) {
	return c.getContractState(id)
}

// getContractState is an internal representation of GetContractStateBy* methods.
func (c *Client) getContractState(param any) (*state.Contract, error) {
	var (
		params = []any{param}
		resp   = &state.Contract{}
	)
	if err := c.performRequest("getcontractstate", params, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetNativeContracts queries information about native contracts.
func (c *Client) GetNativeContracts() ([]state.NativeContract, error) {
	var resp []state.NativeContract
	if err := c.performRequest("getnativecontracts", nil, &resp); err != nil {
		return resp, err
	}

	// Update native contract hashes.
	c.cacheLock.Lock()
	for _, cs := range resp {
		c.cache.nativeHashes[cs.Manifest.Name] = cs.Hash
	}
	c.cacheLock.Unlock()

	return resp, nil
}

// GetNEP11Balances is a wrapper for getnep11balances RPC.
func (c *Client) GetNEP11Balances(address util.Uint160) (*result.NEP11Balances, error) {
	params := []any{address.StringLE()}
	resp := new(result.NEP11Balances)
	if err := c.performRequest("getnep11balances", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetNEP17Balances is a wrapper for getnep17balances RPC.
func (c *Client) GetNEP17Balances(address util.Uint160) (*result.NEP17Balances, error) {
	params := []any{address.StringLE()}
	resp := new(result.NEP17Balances)
	if err := c.performRequest("getnep17balances", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetNEP11Properties is a wrapper for getnep11properties RPC. We recommend using
// nep11 package and Properties method there to receive proper VM types and work with them.
// This method is provided mostly for the sake of completeness. For well-known
// attributes like "description", "image", "name" and "tokenURI" it returns strings,
// while for all others []byte (which can be nil).
func (c *Client) GetNEP11Properties(asset util.Uint160, token []byte) (map[string]any, error) {
	params := []any{asset.StringLE(), hex.EncodeToString(token)}
	resp := make(map[string]any)
	if err := c.performRequest("getnep11properties", params, &resp); err != nil {
		return nil, err
	}
	for k, v := range resp {
		if v == nil {
			continue
		}
		str, ok := v.(string)
		if !ok {
			return nil, errors.New("value is not a string")
		}
		if result.KnownNEP11Properties[k] {
			continue
		}
		val, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, err
		}
		resp[k] = val
	}
	return resp, nil
}

// GetNEP11Transfers is a wrapper for getnep11transfers RPC. Address parameter
// is mandatory, while all others are optional. Limit and page parameters are
// only supported by NeoGo servers and can only be specified with start and stop.
func (c *Client) GetNEP11Transfers(address util.Uint160, start, stop *uint64, limit, page *int) (*result.NEP11Transfers, error) {
	params, err := packTransfersParams(address, start, stop, limit, page)
	if err != nil {
		return nil, err
	}
	resp := new(result.NEP11Transfers)
	if err := c.performRequest("getnep11transfers", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func packTransfersParams(address util.Uint160, start, stop *uint64, limit, page *int) ([]any, error) {
	params := []any{address.StringLE()}
	if start != nil {
		params = append(params, *start)
		if stop != nil {
			params = append(params, *stop)
			if limit != nil {
				params = append(params, *limit)
				if page != nil {
					params = append(params, *page)
				}
			} else if page != nil {
				return nil, errors.New("bad parameters")
			}
		} else if limit != nil || page != nil {
			return nil, errors.New("bad parameters")
		}
	} else if stop != nil || limit != nil || page != nil {
		return nil, errors.New("bad parameters")
	}
	return params, nil
}

// GetNEP17Transfers is a wrapper for getnep17transfers RPC. Address parameter
// is mandatory while all the others are optional. Start and stop parameters
// are supported since neo-go 0.77.0 and limit and page since neo-go 0.78.0.
// These parameters are positional in the JSON-RPC call. For example, you can't specify the limit
// without specifying start/stop first.
func (c *Client) GetNEP17Transfers(address util.Uint160, start, stop *uint64, limit, page *int) (*result.NEP17Transfers, error) {
	params, err := packTransfersParams(address, start, stop, limit, page)
	if err != nil {
		return nil, err
	}
	resp := new(result.NEP17Transfers)
	if err := c.performRequest("getnep17transfers", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPeers returns a list of the nodes that the node is currently connected to/disconnected from.
func (c *Client) GetPeers() (*result.GetPeers, error) {
	var resp = &result.GetPeers{}

	if err := c.performRequest("getpeers", nil, resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetRawMemPool returns a list of unconfirmed transactions in the memory.
func (c *Client) GetRawMemPool() ([]util.Uint256, error) {
	var resp = new([]util.Uint256)

	if err := c.performRequest("getrawmempool", nil, resp); err != nil {
		return *resp, err
	}
	return *resp, nil
}

// GetRawTransaction returns a transaction by hash.
func (c *Client) GetRawTransaction(hash util.Uint256) (*transaction.Transaction, error) {
	var (
		params = []any{hash.StringLE()}
		resp   []byte
		err    error
	)
	if err = c.performRequest("getrawtransaction", params, &resp); err != nil {
		return nil, err
	}
	tx, err := transaction.NewTransactionFromBytes(resp)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetRawTransactionVerbose returns a transaction wrapper with additional
// metadata by transaction's hash.
// NOTE: to get transaction.ID and transaction.Size, use t.Hash() and io.GetVarSize(t) respectively.
func (c *Client) GetRawTransactionVerbose(hash util.Uint256) (*result.TransactionOutputRaw, error) {
	var (
		params = []any{hash.StringLE(), 1} // 1 for verbose.
		resp   = &result.TransactionOutputRaw{}
		err    error
	)
	if err = c.performRequest("getrawtransaction", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetProof returns existence proof of storage item state by the given stateroot
// historical contract hash and historical item key.
func (c *Client) GetProof(stateroot util.Uint256, historicalContractHash util.Uint160, historicalKey []byte) (*result.ProofWithKey, error) {
	var (
		params = []any{stateroot.StringLE(), historicalContractHash.StringLE(), historicalKey}
		resp   = &result.ProofWithKey{}
	)
	if err := c.performRequest("getproof", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// VerifyProof returns value by the given stateroot and proof.
func (c *Client) VerifyProof(stateroot util.Uint256, proof *result.ProofWithKey) ([]byte, error) {
	var (
		params = []any{stateroot.StringLE(), proof.String()}
		resp   []byte
	)
	if err := c.performRequest("verifyproof", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetState returns historical contract storage item state by the given stateroot,
// historical contract hash and historical item key.
func (c *Client) GetState(stateroot util.Uint256, historicalContractHash util.Uint160, historicalKey []byte) ([]byte, error) {
	var (
		params = []any{stateroot.StringLE(), historicalContractHash.StringLE(), historicalKey}
		resp   []byte
	)
	if err := c.performRequest("getstate", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// FindStates returns historical contract storage item states by the given stateroot,
// historical contract hash and historical prefix. If `start` path is specified, items
// starting from `start` path are being returned (excluding item located at the start path).
// If `maxCount` specified, the maximum number of items to be returned equals to `maxCount`.
func (c *Client) FindStates(stateroot util.Uint256, historicalContractHash util.Uint160, historicalPrefix []byte,
	start []byte, maxCount *int) (result.FindStates, error) {
	if historicalPrefix == nil {
		historicalPrefix = []byte{}
	}
	var (
		params = []any{stateroot.StringLE(), historicalContractHash.StringLE(), historicalPrefix}
		resp   result.FindStates
	)
	if start == nil && maxCount != nil {
		start = []byte{}
	}
	if start != nil {
		params = append(params, start)
	}
	if maxCount != nil {
		params = append(params, *maxCount)
	}
	if err := c.performRequest("findstates", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetStateRootByHeight returns the state root for the specified height.
func (c *Client) GetStateRootByHeight(height uint32) (*state.MPTRoot, error) {
	return c.getStateRoot(height)
}

// GetStateRootByBlockHash returns the state root for the block with the specified hash.
func (c *Client) GetStateRootByBlockHash(hash util.Uint256) (*state.MPTRoot, error) {
	return c.getStateRoot(hash)
}

func (c *Client) getStateRoot(param any) (*state.MPTRoot, error) {
	var resp = new(state.MPTRoot)
	if err := c.performRequest("getstateroot", []any{param}, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetStateHeight returns the current validated and local node state height.
func (c *Client) GetStateHeight() (*result.StateHeight, error) {
	var resp = new(result.StateHeight)

	if err := c.performRequest("getstateheight", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetStorageByID returns the stored value according to the contract ID and the stored key.
func (c *Client) GetStorageByID(id int32, key []byte) ([]byte, error) {
	return c.getStorage([]any{id, key})
}

// GetStorageByHash returns the stored value according to the contract script hash and the stored key.
func (c *Client) GetStorageByHash(hash util.Uint160, key []byte) ([]byte, error) {
	return c.getStorage([]any{hash.StringLE(), key})
}

func (c *Client) getStorage(params []any) ([]byte, error) {
	var resp []byte
	if err := c.performRequest("getstorage", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetStorageByIDHistoric returns the historical stored value according to the
// contract ID and, stored key and specified stateroot.
func (c *Client) GetStorageByIDHistoric(root util.Uint256, id int32, key []byte) ([]byte, error) {
	return c.getStorageHistoric([]any{root.StringLE(), id, key})
}

// GetStorageByHashHistoric returns the historical stored value according to the
// contract script hash, the stored key and specified stateroot.
func (c *Client) GetStorageByHashHistoric(root util.Uint256, hash util.Uint160, key []byte) ([]byte, error) {
	return c.getStorageHistoric([]any{root.StringLE(), hash.StringLE(), key})
}

func (c *Client) getStorageHistoric(params []any) ([]byte, error) {
	var resp []byte
	if err := c.performRequest("getstoragehistoric", params, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// FindStorageByHash returns contract storage items by the given contract hash and prefix.
// If `start` index is specified, items starting from `start` index are being returned
// (including item located at the start index).
func (c *Client) FindStorageByHash(contractHash util.Uint160, prefix []byte, start *int) (result.FindStorage, error) {
	var params = []any{contractHash.StringLE(), prefix}
	if start != nil {
		params = append(params, *start)
	}
	return c.findStorage(params)
}

// FindStorageByID returns contract storage items by the given contract ID and prefix.
// If `start` index is specified, items starting from `start` index are being returned
// (including item located at the start index).
func (c *Client) FindStorageByID(contractID int32, prefix []byte, start *int) (result.FindStorage, error) {
	var params = []any{contractID, prefix}
	if start != nil {
		params = append(params, *start)
	}
	return c.findStorage(params)
}

func (c *Client) findStorage(params []any) (result.FindStorage, error) {
	var resp result.FindStorage
	if err := c.performRequest("findstorage", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// FindStorageByHashHistoric returns historical contract storage items by the given stateroot,
// historical contract hash and historical prefix. If `start` index is specified, then items
// starting from `start` index are being returned (including item located at the start index).
func (c *Client) FindStorageByHashHistoric(stateroot util.Uint256, historicalContractHash util.Uint160, historicalPrefix []byte,
	start *int) (result.FindStorage, error) {
	if historicalPrefix == nil {
		historicalPrefix = []byte{}
	}
	var params = []any{stateroot.StringLE(), historicalContractHash.StringLE(), historicalPrefix}
	if start != nil {
		params = append(params, start)
	}
	return c.findStorageHistoric(params)
}

// FindStorageByIDHistoric returns historical contract storage items by the given stateroot,
// historical contract ID and historical prefix. If `start` index is specified, then items
// starting from `start` index are being returned (including item located at the start index).
func (c *Client) FindStorageByIDHistoric(stateroot util.Uint256, historicalContractID int32, historicalPrefix []byte,
	start *int) (result.FindStorage, error) {
	if historicalPrefix == nil {
		historicalPrefix = []byte{}
	}
	var params = []any{stateroot.StringLE(), historicalContractID, historicalPrefix}
	if start != nil {
		params = append(params, start)
	}
	return c.findStorageHistoric(params)
}

func (c *Client) findStorageHistoric(params []any) (result.FindStorage, error) {
	var resp result.FindStorage
	if err := c.performRequest("findstoragehistoric", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetTransactionHeight returns the block index where the transaction is found.
func (c *Client) GetTransactionHeight(hash util.Uint256) (uint32, error) {
	var (
		params = []any{hash.StringLE()}
		resp   uint32
	)
	if err := c.performRequest("gettransactionheight", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetUnclaimedGas returns the unclaimed GAS amount for the specified address.
func (c *Client) GetUnclaimedGas(address string) (result.UnclaimedGas, error) {
	var (
		params = []any{address}
		resp   result.UnclaimedGas
	)
	if err := c.performRequest("getunclaimedgas", params, &resp); err != nil {
		return resp, err
	}
	return resp, nil
}

// GetCandidates returns the current list of NEO candidate node with voting data and
// validator status.
func (c *Client) GetCandidates() ([]result.Candidate, error) {
	var resp = new([]result.Candidate)

	if err := c.performRequest("getcandidates", nil, resp); err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetNextBlockValidators returns the current NEO consensus nodes information and voting data.
func (c *Client) GetNextBlockValidators() ([]result.Validator, error) {
	var resp = new([]result.Validator)

	if err := c.performRequest("getnextblockvalidators", nil, resp); err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetVersion returns the version information about the queried node.
func (c *Client) GetVersion() (*result.Version, error) {
	var resp = &result.Version{}

	if err := c.performRequest("getversion", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// InvokeScript returns the result of the given script after running it true the VM.
// NOTE: This is a test invoke and will not affect the blockchain.
func (c *Client) InvokeScript(script []byte, signers []transaction.Signer) (*result.Invoke, error) {
	var p = []any{script}
	return c.invokeSomething("invokescript", p, signers)
}

// InvokeScriptAtHeight returns the result of the given script after running it
// true the VM using the provided chain state retrieved from the specified chain
// height.
// NOTE: This is a test invoke and will not affect the blockchain.
func (c *Client) InvokeScriptAtHeight(height uint32, script []byte, signers []transaction.Signer) (*result.Invoke, error) {
	var p = []any{height, script}
	return c.invokeSomething("invokescripthistoric", p, signers)
}

// InvokeScriptWithState returns the result of the given script after running it
// true the VM using the provided chain state retrieved from the specified
// state root or block hash.
// NOTE: This is a test invoke and will not affect the blockchain.
func (c *Client) InvokeScriptWithState(stateOrBlock util.Uint256, script []byte, signers []transaction.Signer) (*result.Invoke, error) {
	var p = []any{stateOrBlock.StringLE(), script}
	return c.invokeSomething("invokescripthistoric", p, signers)
}

// InvokeFunction returns the results after calling the smart contract scripthash
// with the given operation and parameters.
// NOTE: this is test invoke and will not affect the blockchain.
func (c *Client) InvokeFunction(contract util.Uint160, operation string, params []smartcontract.Parameter, signers []transaction.Signer) (*result.Invoke, error) {
	var p = []any{contract.StringLE(), operation, params}
	return c.invokeSomething("invokefunction", p, signers)
}

// InvokeFunctionAtHeight returns the results after calling the smart contract
// with the given operation and parameters at the given blockchain state
// specified by the blockchain height.
// NOTE: this is test invoke and will not affect the blockchain.
func (c *Client) InvokeFunctionAtHeight(height uint32, contract util.Uint160, operation string, params []smartcontract.Parameter, signers []transaction.Signer) (*result.Invoke, error) {
	var p = []any{height, contract.StringLE(), operation, params}
	return c.invokeSomething("invokefunctionhistoric", p, signers)
}

// InvokeFunctionWithState returns the results after calling the smart contract
// with the given operation and parameters at the given blockchain state defined
// by the specified state root or block hash.
// NOTE: this is test invoke and will not affect the blockchain.
func (c *Client) InvokeFunctionWithState(stateOrBlock util.Uint256, contract util.Uint160, operation string, params []smartcontract.Parameter, signers []transaction.Signer) (*result.Invoke, error) {
	var p = []any{stateOrBlock.StringLE(), contract.StringLE(), operation, params}
	return c.invokeSomething("invokefunctionhistoric", p, signers)
}

// InvokeContractVerify returns the results after calling `verify` method of the smart contract
// with the given parameters under verification trigger type.
// NOTE: this is test invoke and will not affect the blockchain.
func (c *Client) InvokeContractVerify(contract util.Uint160, params []smartcontract.Parameter, signers []transaction.Signer, witnesses ...transaction.Witness) (*result.Invoke, error) {
	var p = []any{contract.StringLE(), params}
	return c.invokeSomething("invokecontractverify", p, signers, witnesses...)
}

// InvokeContractVerifyAtHeight returns the results after calling `verify` method
// of the smart contract with the given parameters under verification trigger type
// at the blockchain state specified by the blockchain height.
// NOTE: this is test invoke and will not affect the blockchain.
func (c *Client) InvokeContractVerifyAtHeight(height uint32, contract util.Uint160, params []smartcontract.Parameter, signers []transaction.Signer, witnesses ...transaction.Witness) (*result.Invoke, error) {
	var p = []any{height, contract.StringLE(), params}
	return c.invokeSomething("invokecontractverifyhistoric", p, signers, witnesses...)
}

// InvokeContractVerifyWithState returns the results after calling `verify` method
// of the smart contract with the given parameters under verification trigger type
// at the blockchain state specified by the state root or block hash.
// NOTE: this is test invoke and will not affect the blockchain.
func (c *Client) InvokeContractVerifyWithState(stateOrBlock util.Uint256, contract util.Uint160, params []smartcontract.Parameter, signers []transaction.Signer, witnesses ...transaction.Witness) (*result.Invoke, error) {
	var p = []any{stateOrBlock.StringLE(), contract.StringLE(), params}
	return c.invokeSomething("invokecontractverifyhistoric", p, signers, witnesses...)
}

// invokeSomething is an inner wrapper for Invoke* functions.
func (c *Client) invokeSomething(method string, p []any, signers []transaction.Signer, witnesses ...transaction.Witness) (*result.Invoke, error) {
	var resp = new(result.Invoke)
	if signers != nil {
		if witnesses == nil {
			p = append(p, signers)
		} else {
			if len(witnesses) != len(signers) {
				return nil, fmt.Errorf("number of witnesses should match number of signers, got %d vs %d", len(witnesses), len(signers))
			}
			signersWithWitnesses := make([]neorpc.SignerWithWitness, len(signers))
			for i := range signersWithWitnesses {
				signersWithWitnesses[i] = neorpc.SignerWithWitness{
					Signer:  signers[i],
					Witness: witnesses[i],
				}
			}
			p = append(p, signersWithWitnesses)
		}
	}
	if err := c.performRequest(method, p, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// SendRawTransaction broadcasts the given transaction to the Neo network.
// It always returns transaction hash, when successful (no error) this is the
// hash returned from server, when not it's a locally calculated rawTX hash.
func (c *Client) SendRawTransaction(rawTX *transaction.Transaction) (util.Uint256, error) {
	var (
		params = []any{rawTX.Bytes()}
		resp   = new(result.RelayResult)
	)
	if err := c.performRequest("sendrawtransaction", params, resp); err != nil {
		return rawTX.Hash(), err
	}
	return resp.Hash, nil
}

// SubmitBlock broadcasts a raw block over the NEO network.
func (c *Client) SubmitBlock(b block.Block) (util.Uint256, error) {
	var (
		params []any
		resp   = new(result.RelayResult)
	)
	buf := io.NewBufBinWriter()
	b.EncodeBinary(buf.BinWriter)
	if err := buf.Err; err != nil {
		return util.Uint256{}, err
	}
	params = []any{buf.Bytes()}

	if err := c.performRequest("submitblock", params, resp); err != nil {
		return util.Uint256{}, err
	}
	return resp.Hash, nil
}

// SubmitRawOracleResponse submits a raw oracle response to the oracle node.
// Raw params are used to avoid excessive marshalling.
func (c *Client) SubmitRawOracleResponse(ps []any) error {
	return c.performRequest("submitoracleresponse", ps, new(result.RelayResult))
}

// SignAndPushInvocationTx signs and pushes the given script as an invocation
// transaction using the given wif to sign it and the given cosigners to cosign it if
// possible. It spends the amount of gas specified. It returns a hash of the
// invocation transaction and an error. If one of the cosigners accounts is
// neither contract-based nor unlocked, an error is returned.
//
// Deprecated: please use actor.Actor API, this method will be removed in future
// versions.
func (c *Client) SignAndPushInvocationTx(script []byte, acc *wallet.Account, sysfee int64, netfee fixedn.Fixed8, cosigners []SignerAccount) (util.Uint256, error) {
	tx, err := c.CreateTxFromScript(script, acc, sysfee, int64(netfee), cosigners)
	if err != nil {
		return util.Uint256{}, fmt.Errorf("failed to create tx: %w", err)
	}
	return c.SignAndPushTx(tx, acc, cosigners)
}

// SignAndPushTx signs the given transaction using the given wif and cosigners and pushes
// it to the chain. It returns a hash of the transaction and an error. If one of
// the cosigners accounts is neither contract-based nor unlocked, an error is
// returned.
//
// Deprecated: please use actor.Actor API, this method will be removed in future
// versions.
func (c *Client) SignAndPushTx(tx *transaction.Transaction, acc *wallet.Account, cosigners []SignerAccount) (util.Uint256, error) {
	var (
		txHash util.Uint256
		err    error
	)
	m, err := c.GetNetwork()
	if err != nil {
		return txHash, fmt.Errorf("failed to sign tx: %w", err)
	}
	if err = acc.SignTx(m, tx); err != nil {
		return txHash, fmt.Errorf("failed to sign tx: %w", err)
	}
	// try to add witnesses for the rest of the signers
	for i, signer := range tx.Signers[1:] {
		var isOk bool
		for _, cosigner := range cosigners {
			if signer.Account == cosigner.Signer.Account {
				err = cosigner.Account.SignTx(m, tx)
				if err != nil { // then account is non-contract-based and locked, but let's provide more detailed error
					if paramNum := len(cosigner.Account.Contract.Parameters); paramNum != 0 && cosigner.Account.Contract.Deployed {
						return txHash, fmt.Errorf("failed to add contract-based witness for signer #%d (%s): "+
							"%d parameters must be provided to construct invocation script", i, address.Uint160ToString(signer.Account), paramNum)
					}
					return txHash, fmt.Errorf("failed to add witness for signer #%d (%s): account should be unlocked to add the signature. "+
						"Store partially-signed transaction and then use 'wallet sign' command to cosign it", i, address.Uint160ToString(signer.Account))
				}
				isOk = true
				break
			}
		}
		if !isOk {
			return txHash, fmt.Errorf("failed to add witness for signer #%d (%s): account wasn't provided", i, address.Uint160ToString(signer.Account))
		}
	}
	txHash = tx.Hash()
	actualHash, err := c.SendRawTransaction(tx)
	if err != nil {
		return txHash, fmt.Errorf("failed to send tx: %w", err)
	}
	if !actualHash.Equals(txHash) {
		return actualHash, fmt.Errorf("sent and actual tx hashes mismatch:\n\tsent: %v\n\tactual: %v", txHash.StringLE(), actualHash.StringLE())
	}
	return txHash, nil
}

// getSigners returns an array of transaction signers and corresponding accounts from
// given sender and cosigners. If cosigners list already contains sender, the sender
// will be placed at the start of the list.
func getSigners(sender *wallet.Account, cosigners []SignerAccount) ([]transaction.Signer, []*wallet.Account, error) {
	var (
		signers  []transaction.Signer
		accounts []*wallet.Account
	)
	from := sender.ScriptHash()
	s := transaction.Signer{
		Account: from,
		Scopes:  transaction.None,
	}
	for _, c := range cosigners {
		if c.Signer.Account == from {
			s = c.Signer
			continue
		}
		signers = append(signers, c.Signer)
		accounts = append(accounts, c.Account)
	}
	signers = append([]transaction.Signer{s}, signers...)
	accounts = append([]*wallet.Account{sender}, accounts...)
	return signers, accounts, nil
}

// SignAndPushP2PNotaryRequest creates and pushes a P2PNotary request constructed from the main
// and fallback transactions using the given wif to sign it. It returns the request and an error.
// Fallback transaction is constructed from the given script using the amount of gas specified.
// For successful fallback transaction validation at least 2*transaction.NotaryServiceFeePerKey
// GAS should be deposited to the Notary contract.
// Main transaction should be constructed by the user. Several rules should be met for
// successful main transaction acceptance:
//  1. Native Notary contract should be a signer of the main transaction.
//  2. Notary signer should have None scope.
//  3. Main transaction should have dummy contract witness for Notary signer.
//  4. Main transaction should have NotaryAssisted attribute with NKeys specified.
//  5. NotaryAssisted attribute and dummy Notary witness (as long as the other incomplete witnesses)
//     should be paid for. Use CalculateNotaryWitness to calculate the amount of network fee to pay
//     for the attribute and Notary witness.
//  6. Main transaction either shouldn't have all witnesses attached (in this case none of them
//     can be multisignature), or it only should have a partial multisignature.
//
// Note: client should be initialized before SignAndPushP2PNotaryRequest call.
//
// Deprecated: please use Actor from the notary subpackage. This method will be
// deleted in future versions.
func (c *Client) SignAndPushP2PNotaryRequest(mainTx *transaction.Transaction, fallbackScript []byte, fallbackSysFee int64, fallbackNetFee int64, fallbackValidFor uint32, acc *wallet.Account) (*payload.P2PNotaryRequest, error) {
	var err error
	notaryHash, err := c.GetNativeContractHash(nativenames.Notary)
	if err != nil {
		return nil, fmt.Errorf("failed to get native Notary hash: %w", err)
	}
	from := acc.ScriptHash()
	signers := []transaction.Signer{{Account: notaryHash}, {Account: from}}
	if fallbackSysFee < 0 {
		result, err := c.InvokeScript(fallbackScript, signers)
		if err != nil {
			return nil, fmt.Errorf("can't add system fee to fallback transaction: %w", err)
		}
		if result.State != "HALT" {
			return nil, fmt.Errorf("can't add system fee to fallback transaction: bad vm state %s due to an error: %s", result.State, result.FaultException)
		}
		fallbackSysFee = result.GasConsumed
	}

	maxNVBDelta, err := c.GetMaxNotValidBeforeDelta()
	if err != nil {
		return nil, fmt.Errorf("failed to get MaxNotValidBeforeDelta")
	}
	if int64(fallbackValidFor) > maxNVBDelta {
		return nil, fmt.Errorf("fallback transaction should be valid for not more than %d blocks", maxNVBDelta)
	}
	fallbackTx := transaction.New(fallbackScript, fallbackSysFee)
	fallbackTx.Signers = signers
	fallbackTx.ValidUntilBlock = mainTx.ValidUntilBlock
	fallbackTx.Attributes = []transaction.Attribute{
		{
			Type:  transaction.NotaryAssistedT,
			Value: &transaction.NotaryAssisted{NKeys: 0},
		},
		{
			Type:  transaction.NotValidBeforeT,
			Value: &transaction.NotValidBefore{Height: fallbackTx.ValidUntilBlock - fallbackValidFor + 1},
		},
		{
			Type:  transaction.ConflictsT,
			Value: &transaction.Conflicts{Hash: mainTx.Hash()},
		},
	}

	fallbackTx.Scripts = []transaction.Witness{
		{
			InvocationScript:   append([]byte{byte(opcode.PUSHDATA1), keys.SignatureLen}, make([]byte, keys.SignatureLen)...),
			VerificationScript: []byte{},
		},
		{
			InvocationScript:   []byte{},
			VerificationScript: acc.GetVerificationScript(),
		},
	}
	fallbackTx.NetworkFee, err = c.CalculateNetworkFee(fallbackTx)
	if err != nil {
		return nil, fmt.Errorf("failed to add network fee: %w", err)
	}
	fallbackTx.NetworkFee += fallbackNetFee
	m, err := c.GetNetwork()
	if err != nil {
		return nil, fmt.Errorf("failed to sign fallback tx: %w", err)
	}
	if err = acc.SignTx(m, fallbackTx); err != nil {
		return nil, fmt.Errorf("failed to sign fallback tx: %w", err)
	}
	fallbackHash := fallbackTx.Hash()
	req := &payload.P2PNotaryRequest{
		MainTransaction:     mainTx,
		FallbackTransaction: fallbackTx,
	}
	req.Witness = transaction.Witness{
		InvocationScript:   append([]byte{byte(opcode.PUSHDATA1), keys.SignatureLen}, acc.SignHashable(m, req)...),
		VerificationScript: acc.GetVerificationScript(),
	}
	actualHash, err := c.SubmitP2PNotaryRequest(req)
	if err != nil {
		return req, fmt.Errorf("failed to submit notary request: %w", err)
	}
	if !actualHash.Equals(fallbackHash) {
		return req, fmt.Errorf("sent and actual fallback tx hashes mismatch:\n\tsent: %v\n\tactual: %v", fallbackHash.StringLE(), actualHash.StringLE())
	}
	return req, nil
}

// CalculateNotaryFee calculates network fee for one dummy Notary witness and NotaryAssisted attribute with NKeys specified.
// The result should be added to the transaction's net fee for successful verification.
//
// Deprecated: NeoGo calculatenetworkfee method handles notary fees as well since 0.99.3, so
// this method is just no longer needed and will be removed in future versions.
func (c *Client) CalculateNotaryFee(nKeys uint8) (int64, error) {
	baseExecFee, err := c.GetExecFeeFactor()
	if err != nil {
		return 0, fmt.Errorf("failed to get BaseExecFeeFactor: %w", err)
	}
	feePerByte, err := c.GetFeePerByte()
	if err != nil {
		return 0, fmt.Errorf("failed to get FeePerByte: %w", err)
	}
	feePerKey, err := c.GetNotaryServiceFeePerKey()
	if err != nil {
		return 0, fmt.Errorf("failed to get NotaryServiceFeePerKey: %w", err)
	}
	return int64((nKeys+1))*feePerKey + // fee for NotaryAssisted attribute
			fee.Opcode(baseExecFee, // Notary node witness
				opcode.PUSHDATA1, opcode.RET, // invocation script
				opcode.PUSH0, opcode.SYSCALL, opcode.RET) + // System.Contract.CallNative
			nativeprices.NotaryVerificationPrice*baseExecFee + // Notary witness verification price
			feePerByte*int64(io.GetVarSize(make([]byte, 66))) + // invocation script per-byte fee
			feePerByte*int64(io.GetVarSize([]byte{})), // verification script per-byte fee
		nil
}

// SubmitP2PNotaryRequest submits given P2PNotaryRequest payload to the RPC node.
func (c *Client) SubmitP2PNotaryRequest(req *payload.P2PNotaryRequest) (util.Uint256, error) {
	var resp = new(result.RelayResult)
	bytes, err := req.Bytes()
	if err != nil {
		return util.Uint256{}, fmt.Errorf("failed to encode request: %w", err)
	}
	params := []any{bytes}
	if err := c.performRequest("submitnotaryrequest", params, resp); err != nil {
		return util.Uint256{}, err
	}
	return resp.Hash, nil
}

// ValidateAddress verifies that the address is a correct NEO address.
// Consider using [address] package instead to do it locally.
func (c *Client) ValidateAddress(address string) error {
	var (
		params = []any{address}
		resp   = &result.ValidateAddress{}
	)

	if err := c.performRequest("validateaddress", params, resp); err != nil {
		return err
	}
	if !resp.IsValid {
		return errors.New("validateaddress returned false")
	}
	return nil
}

// CalculateValidUntilBlock calculates ValidUntilBlock field for tx as
// current blockchain height + number of validators. Number of validators
// is the length of blockchain validators list got from GetNextBlockValidators()
// method. Validators count is being cached and updated every 100 blocks.
//
// Deprecated: please use (*Actor).CalculateValidUntilBlock. This method will be
// removed in future versions.
func (c *Client) CalculateValidUntilBlock() (uint32, error) {
	var (
		result          uint32
		validatorsCount uint32
	)
	blockCount, err := c.GetBlockCount()
	if err != nil {
		return result, fmt.Errorf("can't get block count: %w", err)
	}

	c.cacheLock.RLock()
	if c.cache.calculateValidUntilBlock.expiresAt > blockCount {
		validatorsCount = c.cache.calculateValidUntilBlock.validatorsCount
		c.cacheLock.RUnlock()
	} else {
		c.cacheLock.RUnlock()
		validators, err := c.GetNextBlockValidators()
		if err != nil {
			return result, fmt.Errorf("can't get validators: %w", err)
		}
		validatorsCount = uint32(len(validators))
		c.cacheLock.Lock()
		c.cache.calculateValidUntilBlock = calculateValidUntilBlockCache{
			validatorsCount: validatorsCount,
			expiresAt:       blockCount + cacheTimeout,
		}
		c.cacheLock.Unlock()
	}
	return blockCount + validatorsCount + 1, nil
}

// AddNetworkFee adds network fee for each witness script and optional extra
// network fee to transaction. `accs` is an array signer's accounts.
//
// Deprecated: please use CalculateNetworkFee or actor.Actor. This method will
// be removed in future versions.
func (c *Client) AddNetworkFee(tx *transaction.Transaction, extraFee int64, accs ...*wallet.Account) error {
	if len(tx.Signers) != len(accs) {
		return errors.New("number of signers must match number of scripts")
	}
	size := io.GetVarSize(tx)
	var ef int64
	for i, cosigner := range tx.Signers {
		if accs[i].Contract.Deployed {
			res, err := c.InvokeContractVerify(cosigner.Account, []smartcontract.Parameter{}, tx.Signers)
			if err != nil {
				return fmt.Errorf("failed to invoke verify: %w", err)
			}
			r, err := unwrap.Bool(res, err)
			if err != nil {
				return fmt.Errorf("signer #%d: %w", i, err)
			}
			if !r {
				return fmt.Errorf("signer #%d: `verify` returned `false`", i)
			}
			tx.NetworkFee += res.GasConsumed
			size += io.GetVarSize([]byte{}) * 2 // both scripts are empty
			continue
		}

		if ef == 0 {
			var err error
			ef, err = c.GetExecFeeFactor()
			if err != nil {
				return fmt.Errorf("can't get `ExecFeeFactor`: %w", err)
			}
		}
		netFee, sizeDelta := fee.Calculate(ef, accs[i].Contract.Script)
		tx.NetworkFee += netFee
		size += sizeDelta
	}
	fee, err := c.GetFeePerByte()
	if err != nil {
		return err
	}
	tx.NetworkFee += int64(size)*fee + extraFee
	return nil
}

// GetNetwork returns the network magic of the RPC node the client connected to. It
// requires Init to be done first, otherwise an error is returned.
//
// Deprecated: please use GetVersion (it has the same data in the Protocol section)
// or actor subpackage. This method will be removed in future versions.
func (c *Client) GetNetwork() (netmode.Magic, error) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if !c.cache.initDone {
		return 0, errNetworkNotInitialized
	}
	return c.cache.network, nil
}

// StateRootInHeader returns true if the state root is contained in the block header.
// You should initialize Client cache with Init() before calling StateRootInHeader.
//
// Deprecated: please use GetVersion (it has the same data in the Protocol section).
// This method will be removed in future versions.
func (c *Client) StateRootInHeader() (bool, error) {
	c.cacheLock.RLock()
	defer c.cacheLock.RUnlock()

	if !c.cache.initDone {
		return false, errNetworkNotInitialized
	}
	return c.cache.stateRootInHeader, nil
}

// GetNativeContractHash returns native contract hash by its name.
//
// Deprecated: please use native contract subpackages that have hashes directly
// (gas, management, neo, notary, oracle, policy, rolemgmt) or
// GetContractStateByAddressOrName method that will return hash along with other
// data.
func (c *Client) GetNativeContractHash(name string) (util.Uint160, error) {
	c.cacheLock.RLock()
	hash, ok := c.cache.nativeHashes[name]
	c.cacheLock.RUnlock()
	if ok {
		return hash, nil
	}
	cs, err := c.GetContractStateByAddressOrName(name)
	if err != nil {
		return util.Uint160{}, err
	}
	c.cacheLock.Lock()
	c.cache.nativeHashes[name] = cs.Hash
	c.cacheLock.Unlock()
	return cs.Hash, nil
}

// TraverseIterator returns a set of iterator values (maxItemsCount at max) for
// the specified iterator and session. If result contains no elements, then either
// Iterator has no elements or session was expired and terminated by the server.
// If maxItemsCount is non-positive, then config.DefaultMaxIteratorResultItems
// iterator values will be returned using single `traverseiterator` call.
// Note that iterator session lifetime is restricted by the RPC-server
// configuration and is being reset each time iterator is accessed. If session
// won't be accessed within session expiration time, then it will be terminated
// by the RPC-server automatically.
func (c *Client) TraverseIterator(sessionID, iteratorID uuid.UUID, maxItemsCount int) ([]stackitem.Item, error) {
	if maxItemsCount <= 0 {
		maxItemsCount = config.DefaultMaxIteratorResultItems
	}
	var (
		params = []any{sessionID.String(), iteratorID.String(), maxItemsCount}
		resp   []json.RawMessage
	)
	if err := c.performRequest("traverseiterator", params, &resp); err != nil {
		return nil, err
	}
	result := make([]stackitem.Item, len(resp))
	for i, iBytes := range resp {
		itm, err := stackitem.FromJSONWithTypes(iBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal %d-th iterator value: %w", i, err)
		}
		result[i] = itm
	}

	return result, nil
}

// TerminateSession tries to terminate the specified session and returns `true` iff
// the specified session was found on server.
func (c *Client) TerminateSession(sessionID uuid.UUID) (bool, error) {
	var resp bool
	params := []any{sessionID.String()}
	if err := c.performRequest("terminatesession", params, &resp); err != nil {
		return false, err
	}

	return resp, nil
}

// GetRawNotaryTransaction  returns main or fallback transaction from the
// RPC node's notary request pool.
func (c *Client) GetRawNotaryTransaction(hash util.Uint256) (*transaction.Transaction, error) {
	var (
		params = []any{hash.StringLE()}
		resp   []byte
		err    error
	)
	if err = c.performRequest("getrawnotarytransaction", params, &resp); err != nil {
		return nil, err
	}
	return transaction.NewTransactionFromBytes(resp)
}

// GetRawNotaryTransactionVerbose returns main or fallback transaction from the
// RPC node's notary request pool.
// NOTE: to get transaction.ID and transaction.Size, use t.Hash() and
// io.GetVarSize(t) respectively.
func (c *Client) GetRawNotaryTransactionVerbose(hash util.Uint256) (*transaction.Transaction, error) {
	var (
		params = []any{hash.StringLE(), 1} // 1 for verbose.
		resp   = &transaction.Transaction{}
		err    error
	)
	if err = c.performRequest("getrawnotarytransaction", params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetRawNotaryPool returns hashes of main P2PNotaryRequest transactions that
// are currently in the RPC node's notary request pool with the corresponding
// hashes of fallback transactions.
func (c *Client) GetRawNotaryPool() (*result.RawNotaryPool, error) {
	resp := &result.RawNotaryPool{}
	if err := c.performRequest("getrawnotarypool", nil, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
