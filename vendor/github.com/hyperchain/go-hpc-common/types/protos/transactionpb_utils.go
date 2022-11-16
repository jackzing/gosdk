// Copyright 2016-2022 Hyperchain Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protos

// IsConfigTx returns if this tx is corresponding with a config tx.
func (m *Transaction) IsConfigTx() bool {
	return m.GetTxType() == Transaction_CTX || m.GetTxType() == Transaction_ANCHORTX ||
		m.GetTxType() == Transaction_ANCHORTXAUTO || m.GetTxType() == Transaction_TIMEOUTTX
}
