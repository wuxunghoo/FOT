// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"strings"

	"github.com/33cn/chain33/client/mocks"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/types"
	pb "github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

var (
	g    Grpc
	qapi *mocks.QueueProtocolAPI
)

// Addr is an autogenerated mock type for the Addr type
type Addr struct {
	mock.Mock
}

// Network provides a mock function with given fields:
func (_m *Addr) Network() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// String provides a mock function with given fields:
func (_m *Addr) String() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

func init() {
	//addr := "192.168.1.1"
	//remoteIpWhitelist[addr] = true
	//grpcFuncWhitelist["*"] = true
	cfg := types.NewChain33Config(types.GetDefaultCfgstring())
	Init(cfg)
	qapi = new(mocks.QueueProtocolAPI)
	qapi.On("GetConfig", mock.Anything).Return(cfg)
	g.cli.QueueProtocolAPI = qapi
}

func getOkCtx() context.Context {
	addr := new(Addr)
	addr.On("String").Return("192.168.1.1")

	ctx := context.Background()
	pr := &peer.Peer{
		Addr:     addr,
		AuthInfo: nil,
	}
	ctx = peer.NewContext(ctx, pr)
	return ctx
}

func testSendTransactionOk(t *testing.T) {

	var in *types.Transaction
	reply := &types.Reply{IsOk: true, Msg: nil}
	qapi.On("SendTx", in).Return(reply, nil)

	reply, err := g.SendTransaction(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, true, reply.IsOk, "reply should be ok")
}

func TestGrpc_SendTransactionSync(t *testing.T) {
	var tx types.Transaction
	reply := &types.Reply{IsOk: true, Msg: tx.Hash()}
	mockAPI := new(mocks.QueueProtocolAPI)
	mockAPI.On("SendTx", mock.Anything).Return(reply, nil)
	mockAPI.On("QueryTx", mock.Anything).Return(&types.TransactionDetail{}, nil)

	g := Grpc{}
	g.cli.QueueProtocolAPI = mockAPI
	reply, err := g.SendTransactionSync(getOkCtx(), &tx)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, true, reply.IsOk, "reply should be ok")
	assert.Equal(t, tx.Hash(), reply.Msg)
}

func TestSendTransaction(t *testing.T) {
	testSendTransactionOk(t)
}

func testVersionOK(t *testing.T) {
	reply := &types.VersionInfo{Chain33: "6.0.2"}
	qapi.On("Version").Return(reply, nil)
	data, err := g.Version(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, "6.0.2", data.Chain33, "reply should be ok")
}

func TestVersion(t *testing.T) {
	testVersionOK(t)
}

func testGetMemPoolOK(t *testing.T) {
	var in *types.ReqGetMempool
	qapi.On("GetMempool", in).Return(nil, nil)
	data, err := g.GetMemPool(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func Test_GetMemPool(t *testing.T) {
	testGetMemPoolOK(t)
}

func testGetLastMemPoolOK(t *testing.T) {
	qapi.On("GetLastMempool").Return(nil, nil)
	data, err := g.GetLastMemPool(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetLastMemPool(t *testing.T) {
	testGetLastMemPoolOK(t)
}

func testGetProperFeeOK(t *testing.T) {
	var in *types.ReqProperFee
	qapi.On("GetProperFee", in).Return(&types.ReplyProperFee{ProperFee: 1000000}, nil)
	data, err := g.GetProperFee(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, int64(1000000), data.ProperFee)
}

func TestGetProperFee(t *testing.T) {
	testGetProperFeeOK(t)
}

func testQueryChainError(t *testing.T) {
	var in *pb.ChainExecutor

	qapi.On("QueryChain", in).Return(nil, fmt.Errorf("error")).Once()
	_, err := g.QueryChain(getOkCtx(), in)
	assert.EqualError(t, err, "error", "return error")
}

func testQueryChainOK(t *testing.T) {
	var in *pb.ChainExecutor
	var msg types.Message
	var req types.ReqString
	req.Data = "msg"
	msg = &req
	qapi.On("QueryChain", in).Return(msg, nil).Once()
	data, err := g.QueryChain(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, true, data.IsOk, "reply should be ok")
	var decodemsg types.ReqString
	pb.Decode(data.Msg, &decodemsg)
	assert.Equal(t, req.Data, decodemsg.Data)
}

func TestQueryChain(t *testing.T) {
	testQueryChainError(t)
	testQueryChainOK(t)
}

func testGetPeerInfoOK(t *testing.T) {
	qapi.On("PeerInfo", mock.Anything).Return(nil, nil)
	data, err := g.GetPeerInfo(getOkCtx(), &types.P2PGetPeerReq{})
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetPeerInfo(t *testing.T) {
	testGetPeerInfoOK(t)
}

func testNetInfoOK(t *testing.T) {
	qapi.On("GetNetInfo", mock.Anything).Return(nil, nil)
	data, err := g.NetInfo(getOkCtx(), &types.P2PGetNetInfoReq{})
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestNetInfo(t *testing.T) {
	testNetInfoOK(t)
}

func testGetAccountsOK(t *testing.T) {
	qapi.On("ExecWalletFunc", "wallet", "WalletGetAccountList", mock.Anything).Return(&types.WalletAccounts{}, nil)
	_, err := g.GetAccounts(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestGetAccount(t *testing.T) {
	qapi.On("ExecWalletFunc", "wallet", "WalletGetAccount", mock.Anything).Return(&types.WalletAccount{}, nil)
	_, err := g.GetAccount(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}
func TestGetAccounts(t *testing.T) {
	testGetAccountsOK(t)
}

func testNewAccountOK(t *testing.T) {
	var in *pb.ReqNewAccount
	qapi.On("ExecWalletFunc", "wallet", "NewAccount", in).Return(&types.WalletAccount{}, nil)
	_, err := g.NewAccount(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestNewAccount(t *testing.T) {
	testNewAccountOK(t)
}

func testWalletTransactionListOK(t *testing.T) {
	var in *pb.ReqWalletTransactionList
	qapi.On("ExecWalletFunc", "wallet", "WalletTransactionList", in).Return(&pb.WalletTxDetails{}, nil)
	_, err := g.WalletTransactionList(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestWalletTransactionList(t *testing.T) {
	testWalletTransactionListOK(t)
}

func testImportPrivKeyOK(t *testing.T) {
	var in *pb.ReqWalletImportPrivkey
	qapi.On("ExecWalletFunc", "wallet", "WalletImportPrivkey", in).Return(&pb.WalletAccount{}, nil)
	_, err := g.ImportPrivkey(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestImportPrivKey(t *testing.T) {
	testImportPrivKeyOK(t)
}

func testSendToAddressOK(t *testing.T) {
	var in *pb.ReqWalletSendToAddress
	qapi.On("ExecWalletFunc", "wallet", "WalletSendToAddress", in).Return(&pb.ReplyHash{}, nil)
	_, err := g.SendToAddress(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestSendToAddress(t *testing.T) {
	testSendToAddressOK(t)
}

func testSetTxFeeOK(t *testing.T) {
	var in *pb.ReqWalletSetFee
	qapi.On("ExecWalletFunc", "wallet", "WalletSetFee", in).Return(&pb.Reply{}, nil)
	_, err := g.SetTxFee(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestSetTxFee(t *testing.T) {
	testSetTxFeeOK(t)
}

func testSetLablOK(t *testing.T) {
	var in *pb.ReqWalletSetLabel
	qapi.On("ExecWalletFunc", "wallet", "WalletSetLabel", in).Return(&pb.WalletAccount{}, nil)
	_, err := g.SetLabl(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestSetLabl(t *testing.T) {
	testSetLablOK(t)
}

func testMergeBalanceOK(t *testing.T) {
	var in *pb.ReqWalletMergeBalance
	qapi.On("ExecWalletFunc", "wallet", "WalletMergeBalance", in).Return(&pb.ReplyHashes{}, nil)
	_, err := g.MergeBalance(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestMergeBalance(t *testing.T) {
	testMergeBalanceOK(t)
}

func testSetPasswdOK(t *testing.T) {
	var in *pb.ReqWalletSetPasswd
	qapi.On("ExecWalletFunc", "wallet", "WalletSetPasswd", in).Return(&pb.Reply{}, nil)
	_, err := g.SetPasswd(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestSetPasswd(t *testing.T) {
	testSetPasswdOK(t)
}

func testLockOK(t *testing.T) {
	var in *pb.ReqNil
	qapi.On("ExecWalletFunc", "wallet", "WalletLock", in).Return(&pb.Reply{}, nil)
	_, err := g.Lock(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestLock(t *testing.T) {
	testLockOK(t)
}

func testUnLockOK(t *testing.T) {
	var in *pb.WalletUnLock
	qapi.On("ExecWalletFunc", "wallet", "WalletUnLock", in).Return(&pb.Reply{}, nil)
	_, err := g.UnLock(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestUnLock(t *testing.T) {
	testUnLockOK(t)
}

func testGenSeedOK(t *testing.T) {
	var in *pb.GenSeedLang
	qapi.On("ExecWalletFunc", "wallet", "GenSeed", in).Return(&pb.ReplySeed{}, nil)
	_, err := g.GenSeed(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestGenSeed(t *testing.T) {
	testGenSeedOK(t)
}

func testGetSeedOK(t *testing.T) {
	var in *pb.GetSeedByPw
	qapi.On("ExecWalletFunc", "wallet", "GetSeed", in).Return(&pb.ReplySeed{}, nil)
	_, err := g.GetSeed(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestGetSeed(t *testing.T) {
	testGetSeedOK(t)
}

func testSaveSeedOK(t *testing.T) {
	var in *pb.SaveSeedByPw
	qapi.On("ExecWalletFunc", "wallet", "SaveSeed", in).Return(&pb.Reply{}, nil)
	_, err := g.SaveSeed(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestSaveSeed(t *testing.T) {
	testSaveSeedOK(t)
}

func testGetWalletStatusOK(t *testing.T) {
	var in *pb.ReqNil
	qapi.On("ExecWalletFunc", "wallet", "GetWalletStatus", in).Return(&pb.WalletStatus{}, nil)
	_, err := g.GetWalletStatus(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestGetWalletStatus(t *testing.T) {
	testGetWalletStatusOK(t)
}

func testDumpPrivkeyOK(t *testing.T) {
	var in *pb.ReqString
	qapi.On("ExecWalletFunc", "wallet", "DumpPrivkey", in).Return(&pb.ReplyString{}, nil)
	_, err := g.DumpPrivkey(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestDumpPrivkey(t *testing.T) {
	testDumpPrivkeyOK(t)
}

func testDumpPrivkeysFileOK(t *testing.T) {
	var in *pb.ReqPrivkeysFile
	qapi.On("ExecWalletFunc", "wallet", "DumpPrivkeysFile", in).Return(&pb.Reply{}, nil)
	_, err := g.DumpPrivkeysFile(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestDumpPrivkeysFile(t *testing.T) {
	testDumpPrivkeysFileOK(t)
}

func testImportPrivkeysFileOK(t *testing.T) {
	var in *pb.ReqPrivkeysFile
	qapi.On("ExecWalletFunc", "wallet", "ImportPrivkeysFile", in).Return(&pb.Reply{}, nil)
	_, err := g.ImportPrivkeysFile(getOkCtx(), nil)
	assert.Nil(t, err, "the error should be nil")
}

func TestImportPrivkeysFile(t *testing.T) {
	testImportPrivkeysFileOK(t)
}

func testGetBlocksError(t *testing.T) {
	var in = pb.ReqBlocks{IsDetail: true}
	qapi.On("GetBlocks", &in).Return(nil, fmt.Errorf("error")).Once()
	_, err := g.GetBlocks(getOkCtx(), &in)
	assert.EqualError(t, err, "error", "the error should be error")
}

func testGetBlocksOK(t *testing.T) {
	var in = pb.ReqBlocks{IsDetail: true}
	var details types.BlockDetails

	var block = &types.Block{Version: 1}
	var detail = &types.BlockDetail{Block: block}
	details.Items = append(details.Items, detail)

	qapi.On("GetBlocks", &in).Return(&details, nil).Once()
	data, err := g.GetBlocks(getOkCtx(), &in)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, true, data.IsOk)

	var details2 types.BlockDetails
	pb.Decode(data.Msg, &details2)
	if !proto.Equal(&details, &details2) {
		assert.Equal(t, types.Encode(&details), types.Encode(&details2))
	}
}

func TestGetBlocks(t *testing.T) {
	testGetBlocksError(t)
	testGetBlocksOK(t)
}

func testGetHexTxByHashError(t *testing.T) {
	var in *pb.ReqHash

	qapi.On("QueryTx", in).Return(nil, fmt.Errorf("error")).Once()
	_, err := g.GetHexTxByHash(getOkCtx(), in)
	assert.EqualError(t, err, "error", "the error should be error")
}

func testGetHexTxByHashOK(t *testing.T) {
	var in *pb.ReqHash
	tx := &types.Transaction{Fee: 1}
	var td = &types.TransactionDetail{Tx: tx}
	var tdNil = &types.TransactionDetail{Tx: nil}

	encodetx := common.ToHex(pb.Encode(tx))

	qapi.On("QueryTx", in).Return(tdNil, nil).Once()
	data, err := g.GetHexTxByHash(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, "", data.Tx)

	qapi.On("QueryTx", in).Return(td, nil).Once()
	data, err = g.GetHexTxByHash(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Equal(t, encodetx, data.Tx)
}

func TestGetHexTxByHash(t *testing.T) {
	testGetHexTxByHashError(t)
	testGetHexTxByHashOK(t)
}

func testGetTransactionByAddrOK(t *testing.T) {
	var in *pb.ReqAddr
	qapi.On("GetTransactionByAddr", in).Return(nil, nil)
	data, err := g.GetTransactionByAddr(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetTransactionByAddr(t *testing.T) {
	testGetTransactionByAddrOK(t)
}

func testGetTransactionByHashesOK(t *testing.T) {
	var in *pb.ReqHashes
	qapi.On("GetTransactionByHash", in).Return(nil, nil)
	data, err := g.GetTransactionByHashes(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetTransactionByHashes(t *testing.T) {
	testGetTransactionByHashesOK(t)
}

func testGetHeadersOK(t *testing.T) {
	var in *pb.ReqBlocks
	qapi.On("GetHeaders", in).Return(nil, nil)
	data, err := g.GetHeaders(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetHeaders(t *testing.T) {
	testGetHeadersOK(t)
}

func testGetBlockOverviewOK(t *testing.T) {
	var in *pb.ReqHash
	qapi.On("GetBlockOverview", in).Return(nil, nil)
	data, err := g.GetBlockOverview(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetBlockOverview(t *testing.T) {
	testGetBlockOverviewOK(t)
}

func testGetBlockHashOK(t *testing.T) {
	var in *pb.ReqInt
	qapi.On("GetBlockHash", in).Return(nil, nil)
	data, err := g.GetBlockHash(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetBlockHash(t *testing.T) {
	testGetBlockHashOK(t)
}

func testIsSyncOK(t *testing.T) {
	var in *pb.ReqNil
	qapi.On("IsSync").Return(nil, nil)
	data, err := g.IsSync(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestIsSync(t *testing.T) {
	testIsSyncOK(t)
}

func testIsNtpClockSyncOK(t *testing.T) {
	var in *pb.ReqNil
	qapi.On("IsNtpClockSync").Return(nil, nil)
	data, err := g.IsNtpClockSync(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestIsNtpClockSync(t *testing.T) {
	testIsNtpClockSyncOK(t)
}

func testGetLastHeaderOK(t *testing.T) {
	var in *pb.ReqNil
	qapi.On("GetLastHeader").Return(nil, nil)
	data, err := g.GetLastHeader(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestGetLastHeader(t *testing.T) {
	testGetLastHeaderOK(t)
}

func testQueryTransactionOK(t *testing.T) {
	var in *pb.ReqHash
	qapi.On("QueryTx", in).Return(nil, nil)
	data, err := g.QueryTransaction(getOkCtx(), in)
	assert.Nil(t, err, "the error should be nil")
	assert.Nil(t, data)
}

func TestQueryTransaction(t *testing.T) {
	testQueryTransactionOK(t)
}

func TestReWriteRawTx(t *testing.T) {
	txHex1 := "0a05636f696e73122c18010a281080c2d72f222131477444795771577233553637656a7663776d333867396e7a6e7a434b58434b7120a08d0630a696c0b3f78dd9ec083a2131477444795771577233553637656a7663776d333867396e7a6e7a434b58434b71"
	in := &types.ReWriteRawTx{
		Tx:     txHex1,
		Fee:    29977777777,
		Expire: "130s",
		To:     "aabbccdd",
		Index:  0,
	}

	data, err := g.ReWriteRawTx(getOkCtx(), in)
	assert.Nil(t, err)
	assert.NotNil(t, data.Data)
	rtTx := hex.EncodeToString(data.Data)
	assert.NotEqual(t, txHex1, rtTx)

	tx := &types.Transaction{}
	err = types.Decode(data.Data, tx)
	assert.Nil(t, err)
	assert.Equal(t, tx.Fee, in.Fee)
	assert.Equal(t, in.To, tx.To)
}

func TestGrpc_CreateNoBalanceTransaction(t *testing.T) {
	_, err := g.CreateNoBalanceTransaction(getOkCtx(), &pb.NoBalanceTx{})
	assert.NoError(t, err)
}

func TestGrpc_CreateNoBalanceTxs(t *testing.T) {
	_, err := g.CreateNoBalanceTxs(getOkCtx(), &pb.NoBalanceTxs{TxHexs: []string{"0a05746f6b656e12413804223d0a0443434e5910a09c011a0d74657374207472616e73666572222231333559774e715367694551787577586650626d526d48325935334564673864343820a08d0630969a9fe6c4b9c7ba5d3a2231333559774e715367694551787577586650626d526d483259353345646738643438", "0a05746f6b656e12413804223d0a0443434e5910b0ea011a0d74657374207472616e73666572222231333559774e715367694551787577586650626d526d48325935334564673864343820a08d0630bca0a2dbc0f182e06f3a2231333559774e715367694551787577586650626d526d483259353345646738643438"}})
	assert.NoError(t, err)
}

func TestGrpc_CreateRawTransaction(t *testing.T) {
	_, err := g.CreateRawTransaction(getOkCtx(), &pb.CreateTx{})
	assert.NoError(t, err)
}

func TestGrpc_CreateTransaction(t *testing.T) {
	_, err := g.CreateTransaction(getOkCtx(), &pb.CreateTxIn{Execer: []byte("coins")})
	assert.Equal(t, err, types.ErrActionNotSupport)
}

func TestGrpc_CreateRawTxGroup(t *testing.T) {
	_, err := g.CreateRawTxGroup(getOkCtx(), &pb.CreateTransactionGroup{})
	assert.Equal(t, types.ErrTxGroupCountLessThanTwo, err)
}

func TestGrpc_GetAddrOverview(t *testing.T) {
	_, err := g.GetAddrOverview(getOkCtx(), &types.ReqAddr{})
	assert.Equal(t, err, types.ErrInvalidAddress)
}

func TestGrpc_GetBalance(t *testing.T) {
	qapi.On("StoreGet", mock.Anything).Return(nil, types.ErrInvalidParam)
	_, err := g.GetBalance(getOkCtx(), &types.ReqBalance{})
	assert.Equal(t, err, types.ErrInvalidParam)
}

func TestGrpc_GetAllExecBalance(t *testing.T) {
	_, err := g.GetAllExecBalance(getOkCtx(), &pb.ReqAllExecBalance{})
	assert.Equal(t, err, types.ErrInvalidAddress)
}

func TestGrpc_QueryConsensus(t *testing.T) {
	qapi.On("QueryConsensus", mock.Anything).Return(&types.ReqString{Data: "test"}, nil)
	_, err := g.QueryConsensus(getOkCtx(), &pb.ChainExecutor{})
	assert.NoError(t, err)
}

func TestGrpc_ExecWallet(t *testing.T) {
	qapi.On("ExecWallet", mock.Anything).Return(&types.ReqString{Data: "test"}, nil)
	_, err := g.ExecWallet(getOkCtx(), &pb.ChainExecutor{})
	assert.NoError(t, err)
}

func TestGrpc_GetLastBlockSequence(t *testing.T) {
	qapi.On("GetLastBlockSequence", mock.Anything).Return(nil, nil)
	_, err := g.GetLastBlockSequence(getOkCtx(), &types.ReqNil{})
	assert.NoError(t, err)
}

func TestGrpc_GetBlockByHashes(t *testing.T) {
	qapi.On("GetBlockByHashes", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	_, err := g.GetBlockByHashes(getOkCtx(), &types.ReqHashes{})
	assert.NoError(t, err)
}

func TestGrpc_GetSequenceByHash(t *testing.T) {
	qapi.On("GetSequenceByHash", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	_, err := g.GetSequenceByHash(getOkCtx(), &pb.ReqHash{})
	assert.NoError(t, err)
}

func TestGrpc_SignRawTx(t *testing.T) {
	qapi.On("ExecWalletFunc", "wallet", "SignRawTx", mock.Anything).Return(&pb.ReplySignRawTx{}, nil)
	_, err := g.SignRawTx(getOkCtx(), &types.ReqSignRawTx{})
	assert.NoError(t, err)
}

func TestGrpc_QueryRandNum(t *testing.T) {
	qapi.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(&pb.ReplyHash{Hash: []byte("test")}, nil)
	_, err := g.QueryRandNum(getOkCtx(), &pb.ReqRandHash{})
	assert.NoError(t, err)
}

func TestGrpc_GetFork(t *testing.T) {
	types.RegFork("para", func(cfg *types.Chain33Config) {
		cfg.SetDappFork("para", "fork100", 100)
	})

	str := types.GetDefaultCfgstring()
	newstr := strings.Replace(str, "Title=\"local\"", "Title=\"chain33\"", 1)
	cfg := types.NewChain33Config(newstr)
	Init(cfg)
	api := new(mocks.QueueProtocolAPI)
	api.On("GetConfig", mock.Anything).Return(cfg)
	grpc := Grpc{}
	grpc.cli.QueueProtocolAPI = api
	val, err := grpc.GetFork(getOkCtx(), &pb.ReqKey{Key: []byte("para-fork100")})
	assert.NoError(t, err)
	assert.Equal(t, int64(100), val.Data)

	cfg1 := types.NewChain33Config(types.GetDefaultCfgstring())
	Init(cfg1)
	api1 := new(mocks.QueueProtocolAPI)
	api1.On("GetConfig", mock.Anything).Return(cfg1)
	grpc1 := Grpc{}
	grpc1.cli.QueueProtocolAPI = api1
	val, err = grpc1.GetFork(getOkCtx(), &pb.ReqKey{Key: []byte("ForkBlockHash")})
	assert.NoError(t, err)
	assert.Equal(t, int64(1), val.Data)
}

func TestGrpc_LoadParaTxByTitle(t *testing.T) {
	qapi.On("LoadParaTxByTitle", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	_, err := g.LoadParaTxByTitle(getOkCtx(), &pb.ReqHeightByTitle{})
	assert.NoError(t, err)
}

func TestGrpc_GetParaTxByHeight(t *testing.T) {
	qapi.On("GetParaTxByHeight", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	_, err := g.GetParaTxByHeight(getOkCtx(), &pb.ReqParaTxByHeight{})
	assert.NoError(t, err)
}

func TestGrpc_GetServerTime(t *testing.T) {
	_, err := g.GetServerTime(getOkCtx(), nil)
	assert.NoError(t, err)
}

func TestGrpc_GetCryptoList(t *testing.T) {
	qapi.On("GetCryptoList").Return(nil)
	_, err := g.GetCryptoList(getOkCtx(), nil)
	assert.NoError(t, err)
}

func TestGrpc_SendDelayTransaction(t *testing.T) {
	qapi.On("SendDelayTx", mock.Anything, mock.Anything).Return(nil, nil)
	_, err := g.SendDelayTransaction(getOkCtx(), nil)
	assert.NoError(t, err)
}

func TestGrpc_GetChainConfig(t *testing.T) {
	cfg, err := g.GetChainConfig(getOkCtx(), nil)
	assert.NoError(t, err)
	assert.Equal(t, types.DefaultCoinPrecision, cfg.GetCoinPrecision())
}

func TestGrpc_SendTransactions(t *testing.T) {

	txCount := 10
	in := &types.Transactions{Txs: make([]*types.Transaction, txCount)}
	testMsg := []byte("test")
	qapi.On("SendTx", mock.Anything).Return(&types.Reply{IsOk: true, Msg: testMsg}, nil)

	reply, err := g.SendTransactions(getOkCtx(), in)
	require.Nil(t, err)
	require.Equal(t, txCount, len(reply.GetHashes()))
	require.Equal(t, testMsg, reply.GetHashes()[0])
}
