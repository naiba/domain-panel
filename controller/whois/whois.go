package whois

import (
	"log"
	"net/http"
	"time"

	"git.cm/nb/domain-panel/pkg/mygin"

	"git.cm/nb/domain-panel"
	"github.com/gin-gonic/gin"
	whois "github.com/likexian/whois-go"
	parser "github.com/likexian/whois-parser-go"
)

//Whois whois 查询
func Whois(c *gin.Context) {
	domain := c.Param("domain")
	if !panel.DomainRegexp.Match([]byte(domain)) {
		c.String(http.StatusForbidden, "域名格式不符合规范")
		return
	}
	u := c.MustGet(mygin.KUser).(panel.User)
	if u.Expire.Before(time.Now()) {
		c.String(http.StatusForbidden, "会员到期")
		return
	}
	result, err := whois.Whois(domain)
	if err == nil {
		var parsed parser.WhoisInfo
		parsed, err = parser.Parse(result)
		if err == nil {
			c.JSON(http.StatusOK, parsed)
			return
		}
	}
	log.Println("whois", err)
	c.Status(http.StatusNoContent)
}