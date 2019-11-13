package main

import (
	"errors"
	"time"

	whois "github.com/likexian/whois-go"
	parser "github.com/likexian/whois-parser-go"
	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/controller"
	"github.com/naiba/nbdomain/model"
)

func init() {
	nbdomain.DB.AutoMigrate(
		model.User{},
		model.Panel{},
		model.Cat{},
		model.Domain{},
		model.Offer{},
		model.Order{},
	)
}

func main() {
	controller.Web()
	go updateWhois()
	select {}
}

func updateWhois() {
	var domains []model.Domain
	for {
		nbdomain.DB.Where("whois_update is NULL or DATEDIFF(now(),whois_update)>7").Find(&domains)
		for _, domain := range domains {
			result, err := whois.Whois(domain.Domain)
			var create, expire time.Time
			var register string
			if err == nil {
				var parsed parser.WhoisInfo
				parsed, err = parser.Parse(result)
				if err == nil {
					create, _ = parseTime(parsed.Domain.CreatedDate)
					expire, _ = parseTime(parsed.Domain.ExpirationDate)
					register = parsed.Registrar.Name
				}
			}
			nbdomain.DB.Model(&domain).UpdateColumns(model.Domain{
				Registrar:   register,
				Create:      create,
				Expire:      expire,
				WhoisUpdate: time.Now(),
			})
			time.Sleep(time.Minute)
		}
		time.Sleep(time.Hour)
	}
}

var timeLayouts = []string{
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"02.01.2006 15:04:05",
	time.RFC1123,     //= "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    //= "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	time.RFC850,      //= "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC3339,     //= "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, //= "2006-01-02T15:04:05.999999999Z07:00"
}

func parseTime(t string) (tt time.Time, e error) {
	for _, layout := range timeLayouts {
		tt, e = time.Parse(layout, t)
		if e == nil {
			return
		}
	}
	e = errors.New("解析失败")
	return
}