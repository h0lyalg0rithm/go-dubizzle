package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/moovweb/gokogiri"
	"github.com/moovweb/gokogiri/xpath"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
)

// go run dubizzle.go --city ajman search motors --category a  --motor_make toyota

type DubizzleParams struct {
	action, city, category, motor_make string
}

type DubizzleResult struct {
	Title, Price, color, description, Url string
	kilometers                            int
}

func searchAPI(param *string) {
	url := fmt.Sprintf("http://dubai.dubizzle.com/search/?keywords=%s", *param)
	ads, err := getResults(&url)
	if err != nil {
		panic(err)
	}
	displayAd(ads)
}

func getResults(url *string) (*[]DubizzleResult, error) {
	resp, err := getHtml(*url)
	if err != nil {
		return nil, err
	}

	return parseHtml(resp)
}

func getHtml(url string) ([]byte, error) {
	// fetch and read a web page
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	println("got the response: \n")

	page, err := ioutil.ReadAll(resp.Body)
	println("read the page: \n")
	return page, err
}

func parseHtml(page []byte) (*[]DubizzleResult, error) {
	// parse the web page
	doc, err := gokogiri.ParseHtml(page)
	if err != nil {
		return nil, err
	}
	println("parsed the doc: \n")

	// perform operations on the parsed page
	xp := xpath.Compile("//*[@id='results-list']/div")
	result_list, err := doc.Root().Search(xp)
	if err != nil {
		return nil, err
	}

	ads := []DubizzleResult{}
	for _, rslt := range result_list {
		xptitle := xpath.Compile(".//h3[@id='title']/span[@class='title']/a")
		xpprice := xpath.Compile(".//div[@class='price']")
		title_info, title_err := rslt.Search(xptitle)
		price_info, price_err := rslt.Search(xpprice)
		if title_err == nil || price_err == nil {
			if len(title_info) > 0 && len(price_info) > 0 {
				title := title_info[0].InnerHtml()
				price := price_info[0].InnerHtml()
				price = strings.TrimSpace(price)
				price = strings.Trim(price, "<br>")
				price = strings.TrimSpace(price)
				url := title_info[0].Attribute("href").String()
				ads = append(ads, DubizzleResult{Title: title, Price: price, Url: url})
			}
		}
	}
	doc.Free()
	return &ads, err
}

func pingAPI(params DubizzleParams) {
	// Generate the URL for dubizzle
	url := fmt.Sprint("https://", params.city, ".dubizzle.com/", params.action, "/", params.category, "/", params.motor_make, "/")
	println("fetching the data from the URL: ", url)

	// important -- don't forget to free the resources when you're done!
	ads, err := getResults(&url)
	if err != nil {
		panic(err)
	}
	displayAd(ads)
}

func displayAd(ads *[]DubizzleResult) {
	for _, ad := range *ads {
		tmpl, err := template.New("test").Parse("{{.Title}}  | {{.Price}}\n")
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(os.Stdout, ad)
		if err != nil {
			panic(err)
		}
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "go-dubizzle"
	app.Version = "0.0.1"
	app.Usage = "Go and get my data from Dubizzle"

	var city, category, motor_make string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "city",
			Value:       "dubai",
			Usage:       "select the city",
			Destination: &city,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search dubizzle",
			Subcommands: []cli.Command{
				{
					Name:  "motors",
					Usage: "search for motors",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:        "category",
							Value:       "used-cars",
							Usage:       "used cars for sale",
							Destination: &category,
						},
						cli.StringFlag{
							Name:        "motor_make",
							Value:       "lexus",
							Usage:       "motor make",
							Destination: &motor_make,
						},
					},
					Action: func(c *cli.Context) {
						println("Searching Dubizzle for motors. Please wait ....")
						// println("Flags: city:", city, " category", category, "motor_make", motor_make)
						params := DubizzleParams{action: "motors", city: city, category: category, motor_make: motor_make}
						// println(params)
						pingAPI(params)
					},
				},
				{
					Name:  "classifields",
					Usage: "search for classifields",
					Action: func(c *cli.Context) {
						println("Searching Dubizzle for classifields. Please wait ....")
						println("Flags: city:", city)
					},
				},
				{
					Name:  "properties",
					Usage: "search for properties (sale or rent)",
					Action: func(c *cli.Context) {
						println("Searching Dubizzle for properties. Please wait ....")
						println("Flags: city:", city)
					},
				},
				{
					Name:  "jobs",
					Usage: "search for jobs (jobs and jobs wanted)",
					Action: func(c *cli.Context) {
						println("Searching Dubizzle for jobs. Please wait ....")
						println("Flags: city:", city)
					},
				},
			},
		},
		{
			Name:    "quick",
			Aliases: []string{"s"},
			Usage:   "Search on Dubizzle",
			Action: func(c *cli.Context) {
				fmt.Print("Searching for ", c.Args().First())
				query := c.Args().First()
				searchAPI(&query)
			},
		},
		{
			Name:    "login",
			Aliases: []string{"l"},
			Usage:   "login to dubizzle",
			Action: func(c *cli.Context) {
				println("loging in")
				println("username")
				println("password")
			},
		},
	}
	app.Run(os.Args)
}

/*

	Usage
	-----

	# Default search will yeild all used lexus cars in dubai
	go run dubizzle.go search motors

	# List all the used toyota cars in dubai
	go run dubizzle.go --city dubai search motors --category used_cars  --motor_make toyota

	# List all the used mercedes cars in dubai
	go run dubizzle.go --city dubai search motors --category used_cars  --motor_make mercedes

	# List all the used B-Class mercedes vehicles in sharjah
	go run dubizzle.go --city sharjah search motors --category used_cars  --motor_make mercedes --model B-Class

	# List all the used B-Class mercedes vehicles in sharjah with minimum prize 50000 AED and maximum 100000 AED
	go run dubizzle.go --city sharjah search motors --category used_cars  --motor_make mercedes --model B-Class --min_price 50000 --max_price 100000


*/

/*

	Sample Dubizzle motor URL:

	https://dubai.dubizzle.com/motors/used-cars/lexus/
	?seller_type=&is_search=1&kilometers__gte=&price__gte=&year__gte=
	&year__lte=2016&kilometers__lte=&keywords=&is_basic_search_widget=1&price__lte=

*/
