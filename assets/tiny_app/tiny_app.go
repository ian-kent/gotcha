package main
import ("github.com/ian-kent/gotcha/app"; "github.com/ian-kent/gotcha/http")
func main() {
	app.Create(nil).Start().Router.Get("/", func(s *http.Session) {
		s.Response.WriteText("Gotcha!")
	})
	<-make(chan int)
}