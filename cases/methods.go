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
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/native/go_abi/node_manager_abi"
	"github.com/ethereum/go-ethereum/contracts/native/governance/node_manager"
	"github.com/ethereum/go-ethereum/contracts/native/utils"
	"github.com/polynetwork/gov-tool/framework"
	"github.com/polynetwork/gov-tool/log"
	"github.com/polynetwork/gov-tool/zion"
	"io/ioutil"
	"math/big"
	"strings"
)

type CreateValidatorParam struct {
	NodeKey         string
	ConsensusPubkey string
	ProposalAddress string
	Commission      *big.Int
	InitStake       *big.Int
	Desc            string
}

func CreateValidator(ctx *framework.FrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/CreateValidator.json")
	if err != nil {
		log.Error(err)
		return false
	}
	createValidatorParam := new(CreateValidatorParam)
	err = json.Unmarshal(data, createValidatorParam)
	if err != nil {
		log.Error(err)
		return false
	}
	signer, err := zion.NewZionSigner(createValidatorParam.NodeKey)
	if err != nil {
		log.Error(err)
		return false
	}
	proposalAddress := common.HexToAddress(createValidatorParam.ProposalAddress)

	scmAbi, err := abi.JSON(strings.NewReader(node_manager_abi.INodeManagerABI))
	if err != nil {
		log.Errorf("CreateValidator, abi.JSON error:" + err.Error())
		return false
	}
	txData, err := scmAbi.Pack("createValidator", createValidatorParam.ConsensusPubkey, proposalAddress,
		createValidatorParam.Commission, createValidatorParam.InitStake, createValidatorParam.Desc)
	if err != nil {
		log.Errorf("CreateValidator, scmAbi.Pack error:" + err.Error())
		return false
	}

	err = CallZionNative(ctx.Z, signer, utils.NodeManagerContractAddress, txData, 1)
	if err != nil {
		log.Errorf("CreateValidator, scmAbi.Pack error:" + err.Error())
		return false
	}
	return true
}

type ChangeEpochParam struct {
	NodeKey         string
}

func ChangeEpoch(ctx *framework.FrameworkContext) bool {
	data, err := ioutil.ReadFile("./params/ChangeEpoch.json")
	if err != nil {
		log.Error(err)
		return false
	}
	changeEpochParam := new(ChangeEpochParam)
	err = json.Unmarshal(data, changeEpochParam)
	if err != nil {
		log.Error(err)
		return false
	}
	signer, err := zion.NewZionSigner(changeEpochParam.NodeKey)
	if err != nil {
		log.Error(err)
		return false
	}

	param := node_manager.ChangeEpochParam{}
	txData, err := param.Encode()
	err = CallZionNative(ctx.Z, signer, utils.NodeManagerContractAddress, txData, 1)
	if err != nil {
		log.Errorf("EpochChange, scmAbi.Pack error:" + err.Error())
		return false
	}
	return true
}
