// Copyright 2019 xgfone
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

package gconf

import (
	"os"
	"os/signal"
)

// ReloadConfigBySignal watches the signal and reload the config
// by calling these parsers.
func ReloadConfigBySignal(sig os.Signal, conf *Config, parsers ...Parser) {
	if len(parsers) == 0 {
		return
	}

	ss := make(chan os.Signal, 1)
	signal.Notify(ss, sig)
	for {
		<-ss
		for _, parser := range parsers {
			conf.Debugf("[HotReload] Calling the parser '%s'", parser.Name())
			if err := parser.Parse(conf); err != nil {
				conf.Printf("[HotReload] the parser '%s' failed to reload: %v",
					parser.Name(), err)
			}
		}
	}
}
