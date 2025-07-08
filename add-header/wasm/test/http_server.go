package main

import (
	"fmt"
	"log"
	"net/http"
)

var res = `
{
  "results": [
    {
      "app_id": "YTRm5Nzjx51vM5MHsUqKApFGEC9MeoJt",
      "pro_type": "pro",
    },
    {
      "app_id": "36H2B7bOIk3wqDKdN7BGtl0gvr8glAd2",
      "pro_type": "pro",
      }
    }
  ]
}
`

func listProApps(w http.ResponseWriter, r *http.Request) {
	fmt.Println("incomming")
	fmt.Fprintf(w, res)
}

func main() {
	http.HandleFunc("/Professional/list-pro-apps", listProApps)

	err := http.ListenAndServe(":9070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
