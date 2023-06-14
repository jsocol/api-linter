// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/googleapis/api-linter/lint"
)

// formatGitHubActionOutput returns lint errors in GitHub actions format.
func formatGitHubActionOutput(responses []lint.Response) []byte {
	var buf bytes.Buffer
	for _, response := range responses {
		for _, problem := range response.Problems {
			// lint example:
			// ::error file={name},line={line},endLine={endLine},title={title}::{message}
			// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-error-message

			fmt.Fprintf(&buf, "::error file=%s", response.FilePath)
			var location lint.FileLocation
			if problem.Location != nil {
				location = lint.FileLocationFromPBLocation(problem.Location, problem.Descriptor)
				fmt.Fprintf(&buf, ",line=%d,col=%d,endLine=%d,endColumn=%d", location.Start.Line, location.Start.Column, location.End.Line, location.End.Column)
			}

			// GitHub uses :: as control characters (which are also used to delimit
			// linter rules. In order to prevent confusion, replace the double colon
			// with two Armenian full stops which are indistinguishable to my eye.
			runeThatLooksLikeTwoColonsButIsActuallyTwoArmenianFullStops := "%3A%3A"
			title := strings.ReplaceAll(string(problem.RuleID), "::", runeThatLooksLikeTwoColonsButIsActuallyTwoArmenianFullStops)
			message := strings.ReplaceAll(problem.Message, "\n", "%0A")
			uri := problem.GetRuleURI()
			if uri != "" {
				message += "%0A" + uri
			}
			fmt.Fprintf(&buf, ",title=%s::%s", title, message)
			if problem.Location != nil {
				fmt.Fprint(&buf, fmt.Sprintf("<br>line=%d, col=%d, endLine=%d, endCol=%d\n", location.Start.Line, location.Start.Column, location.End.Line, location.End.Column))
			} else {
				fmt.Fprintln(&buf, "%0Ano location info found")
			}
		}
	}

	return buf.Bytes()
}
