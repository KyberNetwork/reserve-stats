package http

import (
	"fmt"

	"github.com/KyberNetwork/reserve-stats/gateway/permission"
	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	scas "github.com/qiangmzsx/string-adapter"
)

//newPermissioner: currently there  is only 2 permission for POST/GET requests
func newPermissioner(readKeyID, writeKeyID string) (gin.HandlerFunc, error) {
	const (
		conf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _ , _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub)  && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
	)
	pol := fmt.Sprintf(`
p, %s, /*, (GET)|(POST)
p, %s, /*, GET
`, writeKeyID, readKeyID)
	sa := scas.NewAdapter(pol)
	e := casbin.NewEnforcer(casbin.NewModel(conf), sa)
	e.LoadPolicy()

	p := permission.NewPermissioner(e)
	return p, nil
}
