/* Package api contains base API implementation of unified alerting
 *
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 *
 * Need to remove unused imports.
 */
package api

import (
	"net/http"

	"github.com/go-macaron/binding"
	apimodels "github.com/grafana/alerting-api/pkg/api"
	"github.com/grafana/grafana/pkg/api/response"
	"github.com/grafana/grafana/pkg/api/routing"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/models"
)

type PermissionsApiService interface {
	RouteGetNamespacePermissions(*models.ReqContext) response.Response
	RouteSetNamespacePermissions(*models.ReqContext, apimodels.Permissions) response.Response
}

type PermissionsApiBase struct {
	log log.Logger
}

func (api *API) RegisterPermissionsApiEndpoints(srv PermissionsApiBase) {
	api.RouteRegister.Group("", func(group routing.RouteRegister) {
		group.Get(toMacaronPath("/api/v1/namespace/{Namespace}/permissions"), routing.Wrap(srv.RouteGetNamespacePermissions))
		group.Post(toMacaronPath("/api/v1/namespace/{Namespace}/permissions"), binding.Bind(apimodels.Permissions{}), routing.Wrap(srv.RouteSetNamespacePermissions))
	})
}

func (base PermissionsApiBase) RouteGetNamespacePermissions(c *models.ReqContext) response.Response {
	namespace := c.Params(":Namespace")
	base.log.Info("RouteGetNamespacePermissions: ", "Namespace", namespace)
	return response.Error(http.StatusNotImplemented, "", nil)
}

func (base PermissionsApiBase) RouteSetNamespacePermissions(c *models.ReqContext, body apimodels.Permissions) response.Response {
	namespace := c.Params(":Namespace")
	base.log.Info("RouteSetNamespacePermissions: ", "Namespace", namespace)
	return response.Error(http.StatusNotImplemented, "", nil)
}
