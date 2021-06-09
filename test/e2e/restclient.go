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

package e2e

import (
	"github.com/go-resty/resty/v2"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

var (
	urlGetNamespaces     = "/namespaces"
	urlBroadcastDatatype = "/namespaces/default/broadcast/datatype"
	urlGetData           = "/namespaces/default/data"
)

func GetNamespaces(client *resty.Client) (*resty.Response, error) {
	return client.R().
		SetResult(&[]fftypes.Namespace{}).
		Get(urlGetNamespaces)
}

func BroadcastDatatype(client *resty.Client, name string) (*resty.Response, error) {
	return client.R().
		SetBody(fftypes.Datatype{
			Name:    name,
			Version: "1",
			Value: fftypes.Byteable(`
			{
				"type": "object",
				"properties": {
					"property1": {
						"type": "string"
					}
				}// Copyright © 2021 Kaleido, Inc.
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
				
				
			}`),
		}).Post(urlBroadcastDatatype)
}

func GetData(client *resty.Client) (*resty.Response, error) {
	return client.R().
		SetResult(&[]fftypes.Data{}).
		Get(urlGetData)
}