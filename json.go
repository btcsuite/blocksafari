// Copyright (c) 2013-2014 Conformal Systems LLC.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"github.com/conformal/btcjson"
)

type Tx struct {
	hash      string     // tx hash
	inputs    int        // inputs (vin)
	poutidx   []string   // previous output index
	scriptsig []string   // scriptsig / coinbase
	outputs   int        // outputs (vout)
	value     []float64  // amount
	spubkey   []string   // script pubkey
	typ       []string   // type
	addresses [][]string // addresses
}

func getBlock(block string, withTx bool) (btcjson.BlockResult, error) {
	var result btcjson.BlockResult

	cmd, err := btcjson.NewGetBlockCmd("blocksafari", block, true, withTx)
	if err != nil {
		return result, err
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return result, err
	}

	reply, err := btcjson.TlsRpcCommand(cfg.RPCUser, cfg.RPCPassword, cfg.RPCServer, msg, pem, false)
	if err != nil {
		return result, err
	}

	if reply.Error != nil {
		return result, reply.Error
	}

	return reply.Result.(btcjson.BlockResult), nil
}

func getBlockCount() (float64, error) {
	cmd, err := btcjson.NewGetBlockCountCmd("blocksafari")
	if err != nil {
		return -1, err
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return -1, err
	}

	reply, err := btcjson.TlsRpcCommand(cfg.RPCUser, cfg.RPCPassword, cfg.RPCServer, msg, pem, false)
	if err != nil {
		return -1, err
	}

	if reply.Error != nil {
		return -1, reply.Error
	}

	return reply.Result.(float64), nil
}

func getBlockHash(idx int64) (string, error) {
	cmd, err := btcjson.NewGetBlockHashCmd("blocksafari", idx)
	if err != nil {
		return "", err
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return "", err
	}

	reply, err := btcjson.TlsRpcCommand(cfg.RPCUser, cfg.RPCPassword, cfg.RPCServer, msg, pem, false)
	if err != nil {
		return "", err
	}

	if reply.Error != nil {
		return "", reply.Error
	}

	return reply.Result.(string), nil
}

func getRawBlock(block string) (interface{}, error) {
	cmd, err := btcjson.NewGetBlockCmd("blocksafari", block)
	if err != nil {
		return nil, err
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	reply, err := btcjson.TlsRpcCommand(cfg.RPCUser, cfg.RPCPassword, cfg.RPCServer, msg, pem, false)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		return nil, reply.Error
	}

	return reply.Result, nil
}

func getTx(tx string) (btcjson.TxRawResult, error) {
	var result btcjson.TxRawResult

	cmd, err := btcjson.NewGetRawTransactionCmd("blocksafari", tx, true)
	if err != nil {
		return result, err
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return result, err
	}

	reply, err := btcjson.TlsRpcCommand(cfg.RPCUser, cfg.RPCPassword, cfg.RPCServer, msg, pem, false)
	if err != nil {
		return result, err
	}

	if reply.Error != nil {
		return result, reply.Error
	}

	return reply.Result.(btcjson.TxRawResult), nil
}

func getRawTx(tx string) (interface{}, error) {
	cmd, err := btcjson.NewGetRawTransactionCmd("blocksafari", tx, true)
	if err != nil {
		return nil, err
	}

	msg, err := json.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	reply, err := btcjson.TlsRpcCommand(cfg.RPCUser, cfg.RPCPassword, cfg.RPCServer, msg, pem, false)
	if err != nil {
		return nil, err
	}

	if reply.Error != nil {
		return nil, reply.Error
	}

	return reply.Result, nil
}
