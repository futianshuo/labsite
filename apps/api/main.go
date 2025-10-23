package main

func main() {
	db := mustDB()
	r := initRouter(db)
	_ = r.Run(":9000")
}
