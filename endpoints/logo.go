package endpoints

import (
	"fmt"
	"net/http"
	"os"
)

//GetLogo fetches the Logo from the Backend and serves it
func GetLogo(w http.ResponseWriter, r *http.Request) {
	seperator := string(os.PathSeparator)
	filepath, _ := os.Getwd()

	userPath := fmt.Sprint(filepath, seperator, "userData", seperator, "logo", seperator)
	for _, k := range []string{"svg", "png"} {
		_, err := os.Stat(userPath + "logo" + "." + k)
		if err == nil {
			http.ServeFile(w, r, userPath+"logo."+k)
			break
		}
	}
	fmt.Println(filepath + seperator + "userData" + seperator + "sampleData" + seperator + "logo.svg")
	http.ServeFile(w, r, filepath+seperator+"userData"+seperator+"sampleData"+seperator+"logo.svg")
}

//GetFavicon fetches the Logo from the Backend and serves it
func GetFavicon(w http.ResponseWriter, r *http.Request) {
	seperator := string(os.PathSeparator)
	filepath, _ := os.Getwd()

	userPath := fmt.Sprint(filepath, seperator, "userData", seperator, "logo", seperator)
	for _, k := range []string{"svg", "ico", "png"} {
		_, err := os.Stat(userPath + "logo" + "." + k)
		if err == nil {
			http.ServeFile(w, r, userPath+"logo."+k)
			break
		}
	}
	http.ServeFile(w, r, filepath+seperator+"userData"+seperator+"sampleData"+seperator+"logo.ico")
}
