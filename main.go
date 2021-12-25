package main

import "encoding/json"
import "log"
import "net"
import "os"

// Required modules
import "github.com/DavidGamba/go-getoptions"
import "github.com/fasthttp/router"
import "github.com/oschwald/geoip2-golang"
import "github.com/valyala/fasthttp"

// Debugging
//import "github.com/davecgh/go-spew/spew"

// Custom type
type GeoData struct {
	Country		string		`json:"country"`
	Latitude	float64		`json:"latitude"`
	Longitude	float64		`json:"logitude"`
	Continent	string		`json:"continent"`
	Timezone	string		`json:"timezone"`
	// Optional fields (omitempty)
	StateProv	string		`json:"stateprov,omitempty"`
	StateProvCode	string		`json:"stateprovCode,omitempty"`
	City		string		`json:"city,omitempty"`
}

// Global variables
var geodb *geoip2.Reader

func main() {
  // Parse options
  opt := getoptions.New()

  var dbfile string
  opt.StringVar(&dbfile, "db", "", opt.Required())
  var listen string
  opt.StringVar(&listen, "listen", "127.0.0.1:8000")

  remaining, err := opt.Parse(os.Args[1:])
  if len(os.Args[1:]) == 0 { log.Fatal(opt.Help()) }
  if err != nil { log.Fatalf("Could not parse options: %s\n", err) }
  if len(remaining) > 0 { log.Fatalf("Unsupported parameter: %s\n", remaining) }

  var openerr error
  geodb, openerr = geoip2.Open(dbfile)
  if openerr != nil {
    os.Exit(1)
  }
  defer geodb.Close()

  routes := router.New()
  routes.GET("/{ipaddr}", http_handler_get_ipaddr)

  laserr := fasthttp.ListenAndServe(listen, routes.Handler)
  if laserr != nil {
    log.Fatal(laserr)
  }
}

func get_ip_data (ip string) (GeoData, error) {
  var response GeoData
  record, lookuperr := geodb.City(net.ParseIP(ip))
  if lookuperr != nil {
    return response, lookuperr
  }

  response.Country = record.Country.IsoCode
  response.Latitude = record.Location.Latitude
  response.Longitude = record.Location.Longitude
  response.Continent = record.Continent.Code
  response.Timezone = record.Location.TimeZone
  // Optional fields
  if _, ok := record.City.Names["en"]; ok {
    response.City = record.City.Names["en"]
  }
  if len(record.Subdivisions) >= 1 {
    if _, ok := record.Subdivisions[0].Names["en"]; ok {
      response.StateProv = record.Subdivisions[0].Names["en"]
    }
    response.StateProvCode = record.Subdivisions[0].IsoCode
  }

  return response, nil
}

func http_handler_get_ipaddr (ctx *fasthttp.RequestCtx) {
  geodata, geoerr := get_ip_data(ctx.UserValue("ipaddr").(string))
  if geoerr != nil {
	  ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	  ctx.Response.Header.Set("X-Error", geoerr.Error())
	  _, werr := ctx.Write([]byte(`{}`))
	  if werr != nil { log.Println(werr.Error()) }
	  return
  }

  jsondata, jsonerr := json.MarshalIndent(geodata, "", "  ")
  if jsonerr != nil {
	  ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	  return
  }

  _, werr := ctx.Write(jsondata)
  if werr != nil { log.Println(werr.Error()) }
}
