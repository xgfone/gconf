// Copyright 2021 xgfone
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

// Package gconf is an extensible and powerful go configuration manager.
//
// Features
//
//   - A atomic key-value configuration center.
//   - Support kinds of decoders to decode the data from the source.
//   - Support to get the configuration data from many data sources.
//   - Support to change of the configuration option thread-safely during running.
//   - Support to observe the change of the configration options.
//
package gconf
