/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package cases

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/polynetwork/gov-tool/log"
	"github.com/polynetwork/gov-tool/zion"
	"math/big"
	"time"
)

func CallZionNative(z *zion.ZionTools, signer *zion.ZionSigner, contractAddress common.Address, txData []byte, gasPriceMul uint64) error {
	gasPrice, err := z.GetEthClient().SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("CallZionNative, get suggest gas price failed error: %s", err.Error())
	}
	gasPrice = gasPrice.Mul(gasPrice, new(big.Int).SetUint64(gasPriceMul))

	callMsg := ethereum.CallMsg{
		From: signer.Address, To: &contractAddress, Gas: 0, GasPrice: gasPrice,
		Value: big.NewInt(int64(0)), Data: txData,
	}
	gasLimit, err := z.GetEthClient().EstimateGas(context.Background(), callMsg)
	if err != nil {
		return fmt.Errorf("CallZionNative, estimate gas limit error: %s", err.Error())
	}
	nonce := zion.NewNonceManager(z.GetEthClient()).GetAddressNonce(signer.Address)
	tx := types.NewTx(&types.LegacyTx{Nonce: nonce, GasPrice: gasPrice, Gas: gasLimit, To: &contractAddress, Value: big.NewInt(0), Data: txData})
	chainID, err := z.GetChainID()
	if err != nil {
		return fmt.Errorf("CallZionNative, get chain id error: %s", err.Error())
	}
	signedtx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), signer.PrivateKey)
	if err != nil {
		return fmt.Errorf("CallZionNative, SignTransaction failed:%v", err)
	}
	duration := time.Second * 20
	timerCtx, cancelFunc := context.WithTimeout(context.Background(), duration)
	defer cancelFunc()
	err = z.GetEthClient().SendTransaction(timerCtx, signedtx)
	if err != nil {
		return fmt.Errorf("CallZionNative, SendTransaction failed:%v", err)
	}
	txhash := signedtx.Hash()
	isSuccess := z.WaitTransactionConfirm(txhash)
	if isSuccess {
		log.Infof("CallZionNative, success hash: %s", txhash.String())
		return nil
	} else {
		return fmt.Errorf("CallZionNative, failed hash: %s", txhash.String())
	}
}
