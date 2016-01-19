package main

import (
  "fmt"
  "io/ioutil"

  "os"
  "github.com/codegangsta/cli"

  "net/http"
  "github.com/moovweb/gokogiri"
  "github.com/moovweb/gokogiri/xpath"

  "log"
)

// go run dubizzle.go --city ajman search motors --category a  --motor_make toyota

type DubizzleParams struct {
	action, city, category, motor_make string
}

type DubizzleResult struct {
	title, color, description, url string
	price float64
	kilometers int
}

func pingAPI(params DubizzleParams){
	// Generate the URL for dubizzle
	url := fmt.Sprint("https://", params.city, ".dubizzle.com/", params.action, "/", params.category, "/", params.motor_make, "/")
	println("fetching the data from the URL: ", url)

	// fetch and read a web page
  resp, err := http.Get(url)
  if err != nil {
  	log.Fatal(err)
		println("error while fetching: ", err)
		// fmt.Printf("error while fetching: %s\n", string(err))
		os.Exit(1)
	}
	println("got the response: \n")
	defer resp.Body.Close()

  page, err := ioutil.ReadAll(resp.Body)
  if err != nil {
		println("error while reading: ", err)
		// fmt.Printf("error while reading: %s\n", string(err))
		os.Exit(1)
	}
  println("read the page: \n")
  // fmt.Printf("%s\n", string(page))

  // parse the web page
  doc, err := gokogiri.ParseHtml(page)
  if err != nil {
		println("error while parsing: ", err)
		// fmt.Printf("error while parsing: %s\n", string(err))
		os.Exit(1)
	}
	fmt.Printf("parsed the doc: \n")
	
  // perform operations on the parsed page
  xp := xpath.Compile("//*[@id='results-list']/div")
  result_list, err := doc.Root().Search(xp)
  if err != nil {
    log.Println(err)
  }
	for _, rslt := range result_list {
		xptitle := xpath.Compile(".//h3[@id='title']/span[@class='title']/a")
		title_info, err := rslt.Search(xptitle)
		if err != nil {
	    log.Println(err)
	  }
	  if len(title_info) > 0{
	  	fmt.Println(title_info[0].InnerHtml())
	  }
	  
	}

  // important -- don't forget to free the resources when you're done!
  doc.Free()

}

func main() {
	app := cli.NewApp()
	app.Name = "go-dubizzle"
	app.Version = "0.0.1"
	app.Usage = "Go and get my data from Dubizzle"
	
	var city, category, motor_make string
	
	app.Flags = []cli.Flag {
	  cli.StringFlag{
	    Name:        "city",
	    Value:       "dubai",
	    Usage:       "select the city",
	    Destination: &city,
	  },
	}

	app.Commands = []cli.Command{
	  {
	    Name:      "search",
	    Aliases:     []string{"s"},
	    Usage:     "search dubizzle",
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
	    Name:      "login",
	    Aliases:     []string{"l"},
	    Usage:     "login to dubizzle",
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
