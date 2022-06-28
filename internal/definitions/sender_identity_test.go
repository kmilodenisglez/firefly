// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package definitions

import (
	"fmt"
	"testing"

	"github.com/hyperledger/firefly/internal/identity"
	"github.com/hyperledger/firefly/mocks/broadcastmocks"
	"github.com/hyperledger/firefly/mocks/identitymanagermocks"
	"github.com/hyperledger/firefly/mocks/sysmessagingmocks"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDefineIdentity(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()

	mim := ds.identity.(*identitymanagermocks.Manager)
	mbm := ds.broadcast.(*broadcastmocks.Manager)
	mms := &sysmessagingmocks.MessageSender{}

	mim.On("NormalizeSigningKey", mock.Anything, "0x1234", identity.KeyNormalizationBlockchainPlugin).Return("", nil)
	mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("SendAndWait", mock.Anything).Return(nil)

	ds.multiparty = true

	err := ds.DefineIdentity(ds.ctx, &core.IdentityClaim{
		Identity: &core.Identity{},
	}, &core.SignerRef{
		Key: "0x1234",
	}, nil, core.SystemTagDefineNamespace, true)
	assert.NoError(t, err)

	mim.AssertExpectations(t)
	mbm.AssertExpectations(t)
	mms.AssertExpectations(t)
}

func TestDefineIdentityFail(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()

	mim := ds.identity.(*identitymanagermocks.Manager)
	mbm := ds.broadcast.(*broadcastmocks.Manager)
	mms := &sysmessagingmocks.MessageSender{}

	mim.On("NormalizeSigningKey", mock.Anything, "0x1234", identity.KeyNormalizationBlockchainPlugin).Return("", nil)
	mbm.On("NewBroadcast", mock.Anything).Return(mms)
	mms.On("SendAndWait", mock.Anything).Return(fmt.Errorf("pop"))

	ds.multiparty = true

	err := ds.DefineIdentity(ds.ctx, &core.IdentityClaim{
		Identity: &core.Identity{},
	}, &core.SignerRef{
		Key: "0x1234",
	}, nil, core.SystemTagDefineNamespace, true)
	assert.EqualError(t, err, "pop")

	mim.AssertExpectations(t)
	mbm.AssertExpectations(t)
	mms.AssertExpectations(t)
}

func TestDefineIdentityFailKey(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()

	mim := ds.identity.(*identitymanagermocks.Manager)

	mim.On("NormalizeSigningKey", mock.Anything, "0x1234", identity.KeyNormalizationBlockchainPlugin).Return("", fmt.Errorf("pop"))

	ds.multiparty = true

	err := ds.DefineIdentity(ds.ctx, &core.IdentityClaim{
		Identity: &core.Identity{},
	}, &core.SignerRef{
		Key: "0x1234",
	}, nil, core.SystemTagDefineNamespace, true)
	assert.EqualError(t, err, "pop")

	mim.AssertExpectations(t)
}

func TestDefineIdentityChild(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()

	mim := ds.identity.(*identitymanagermocks.Manager)
	mbm := ds.broadcast.(*broadcastmocks.Manager)
	mms1 := &sysmessagingmocks.MessageSender{}
	mms2 := &sysmessagingmocks.MessageSender{}

	mim.On("NormalizeSigningKey", mock.Anything, "0x1234", identity.KeyNormalizationBlockchainPlugin).Return("", nil)
	mbm.On("NewBroadcast", mock.Anything).Return(mms1).Once()
	mbm.On("NewBroadcast", mock.Anything).Return(mms2).Once()
	mms1.On("SendAndWait", mock.Anything).Return(nil)
	mms2.On("Send", mock.Anything).Return(nil)
	mim.On("ResolveInputSigningIdentity", mock.Anything, mock.MatchedBy(func(signer *core.SignerRef) bool {
		return signer.Key == "0x2345"
	})).Return(nil)

	ds.multiparty = true

	err := ds.DefineIdentity(ds.ctx, &core.IdentityClaim{
		Identity: &core.Identity{},
	}, &core.SignerRef{
		Key: "0x1234",
	}, &core.SignerRef{
		Key: "0x2345",
	}, core.SystemTagDefineNamespace, true)
	assert.NoError(t, err)

	mim.AssertExpectations(t)
	mbm.AssertExpectations(t)
	mms1.AssertExpectations(t)
}

func TestDefineIdentityChildFail(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()

	mim := ds.identity.(*identitymanagermocks.Manager)
	mbm := ds.broadcast.(*broadcastmocks.Manager)
	mms1 := &sysmessagingmocks.MessageSender{}
	mms2 := &sysmessagingmocks.MessageSender{}

	mim.On("NormalizeSigningKey", mock.Anything, "0x1234", identity.KeyNormalizationBlockchainPlugin).Return("", nil)
	mbm.On("NewBroadcast", mock.Anything).Return(mms1).Once()
	mbm.On("NewBroadcast", mock.Anything).Return(mms2).Once()
	mms1.On("SendAndWait", mock.Anything).Return(nil)
	mms2.On("Send", mock.Anything).Return(fmt.Errorf("pop"))
	mim.On("ResolveInputSigningIdentity", mock.Anything, mock.MatchedBy(func(signer *core.SignerRef) bool {
		return signer.Key == "0x2345"
	})).Return(nil)

	ds.multiparty = true

	err := ds.DefineIdentity(ds.ctx, &core.IdentityClaim{
		Identity: &core.Identity{},
	}, &core.SignerRef{
		Key: "0x1234",
	}, &core.SignerRef{
		Key: "0x2345",
	}, core.SystemTagDefineNamespace, true)
	assert.EqualError(t, err, "pop")

	mim.AssertExpectations(t)
	mbm.AssertExpectations(t)
	mms1.AssertExpectations(t)
}

func TestDefineIdentityNonMultiparty(t *testing.T) {
	ds, cancel := newTestDefinitionSender(t)
	defer cancel()
	dh := ds.handler

	mim := dh.identity.(*identitymanagermocks.Manager)
	mim.On("VerifyIdentityChain", mock.Anything, mock.AnythingOfType("*core.Identity")).Return(nil, false, fmt.Errorf("pop"))

	ds.multiparty = false

	err := ds.DefineIdentity(ds.ctx, &core.IdentityClaim{
		Identity: &core.Identity{},
	}, &core.SignerRef{
		Key: "0x1234",
	}, nil, core.SystemTagDefineNamespace, true)
	assert.NoError(t, err)

	mim.AssertExpectations(t)
}
