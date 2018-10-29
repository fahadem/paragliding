# paragliding
Assignment 2

#APP heroku 

https://fahadem2.herokuapp.com/

#Problems

MongoDB storage
Clock trigger not implemented but maybe this can work :

package Clock_Trigger

import (
	"net/http"
	"time"
)

func main() {

	for range time.Tick(10 * time.Minute) {
		http.Get("https://fahadem2.herokuapp.com/paragliding/admin/api/webhook")
	}
}

#Urls

Example : 
GET api : https://fahadem2.herokuapp.com/paragliding/api or https://fahadem2.herokuapp.com/paragliding/
 
