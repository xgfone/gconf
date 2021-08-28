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

import "testing"

func TestNewMaybeURLValidator(t *testing.T) {
	validate := NewMaybeURLValidator()
	if err := validate(""); err != nil {
		t.Error(err)
	} else if err = validate("http://www.example.com"); err != nil {
		t.Error(err)
	}
}

func TestNewMaybeIPValidator(t *testing.T) {
	validate := NewMaybeIPValidator()
	if err := validate(""); err != nil {
		t.Error(err)
	} else if err = validate("1.2.3.4"); err != nil {
		t.Error(err)
	}
}

func TestNewMaybeEmailValidator(t *testing.T) {
	validate := NewMaybeEmailValidator()
	if err := validate(""); err != nil {
		t.Error(err)
	} else if err = validate("abc@xyz.com"); err != nil {
		t.Error(err)
	}
}

func TestNewMaybeAddressValidator(t *testing.T) {
	validate := NewMaybeAddressValidator()
	if err := validate(""); err != nil {
		t.Error(err)
	} else if err = validate("1.2.3.4:80"); err != nil {
		t.Error(err)
	}
}

func TestNewAddressOrIPValidator(t *testing.T) {
	validate := NewAddressOrIPValidator()
	if err := validate("1.2.3.4"); err != nil {
		t.Error(err)
	} else if err = validate("1.2.3.4:80"); err != nil {
		t.Error(err)
	}
}

func TestNewMaybeAddressOrIPValidator(t *testing.T) {
	validate := NewMaybeAddressOrIPValidator()
	if err := validate(""); err != nil {
		t.Error(err)
	} else if err = validate("1.2.3.4"); err != nil {
		t.Error(err)
	} else if err = validate("1.2.3.4:80"); err != nil {
		t.Error(err)
	}
}

func TestNewAddressOrIPSliceValidator(t *testing.T) {
	validate := NewAddressOrIPSliceValidator()
	if err := validate([]string{"1.2.3.4", "1.2.3.4:80"}); err != nil {
		t.Error(err)
	}
}
