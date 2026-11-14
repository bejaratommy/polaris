// Copyright 2026 FairwindsOps, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fairwindsops/polaris/pkg/config"
	"github.com/fairwindsops/polaris/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestMainHandlerHidesResourcesWithoutChecks(t *testing.T) {
	success := validator.ResultSet{"check": {Success: true}}
	warning := validator.ResultSet{"check": {Severity: config.SeverityWarning}}
	danger := validator.ResultSet{"check": {Severity: config.SeverityDanger}}
	auditData := validator.AuditData{Results: []validator.Result{
		{Name: "unchecked-config", Namespace: "checked", Kind: "ConfigMap"},
		{Name: "unchecked-deployment", Namespace: "exempt", Kind: "Deployment", PodResult: &validator.PodResult{ContainerResults: []validator.ContainerResult{{Name: "empty"}}}},
		{Name: "unchecked-role", Kind: "ClusterRole"},
		{Name: "checked-spec", Namespace: "checked", Kind: "Deployment", Results: success, PodResult: &validator.PodResult{}},
		{Name: "checked-pod", Namespace: "checked", Kind: "Deployment", PodResult: &validator.PodResult{Results: warning}},
		{Name: "checked-container", Namespace: "checked", Kind: "Deployment", PodResult: &validator.PodResult{ContainerResults: []validator.ContainerResult{{Name: "container", Results: danger}}}},
		{Name: "checked-role", Kind: "ClusterRole", Results: success},
	}}
	originalResults := append([]validator.Result(nil), auditData.Results...)
	for _, query := range []string{"", "?ns=checked", "?ns=exempt"} {
		t.Run(query, func(t *testing.T) {
			response := httptest.NewRecorder()
			MainHandler(response, httptest.NewRequest(http.MethodGet, "/"+query, nil), config.Configuration{}, auditData, "/")
			assert.Equal(t, http.StatusOK, response.Code)
			html := response.Body.String()
			for _, name := range []string{"unchecked-config", "unchecked-deployment", "unchecked-role", "namespace-exempt"} {
				assert.NotContains(t, html, name)
			}
			for _, name := range []string{"checked-spec", "checked-pod", "checked-container"} {
				if query == "?ns=exempt" {
					assert.NotContains(t, html, "<strong>"+name+"</strong>")
				} else {
					assert.Contains(t, html, "<strong>"+name+"</strong>")
				}
			}
			if query == "" {
				assert.Contains(t, html, "<strong>checked-role</strong>")
			} else {
				assert.NotContains(t, html, "<strong>checked-role</strong>")
			}
			assert.Equal(t, originalResults, auditData.Results)
		})
	}
	response := httptest.NewRecorder()
	JSONHandler(response, httptest.NewRequest(http.MethodGet, "/results.json", nil), &auditData)
	assert.Contains(t, response.Body.String(), "unchecked-config")
	assert.Contains(t, response.Body.String(), "unchecked-deployment")
}
