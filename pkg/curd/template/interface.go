/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package template

import (
	"bytes"

	"github.com/hu-1996/cwgo/pkg/curd/code"
)

var interfaceTemplate = `{{.Comment}}
type {{.Name}} interface {
{{.Methods.GetCode}}
}` + "\n"

type InterfaceRender struct {
	Name    string
	Comment string
	Methods code.InterfaceMethods
}

func (ir *InterfaceRender) RenderObj(buffer *bytes.Buffer) error {
	if err := templateRender(buffer, "interfaceTemplate", interfaceTemplate, ir); err != nil {
		return err
	}
	return nil
}
