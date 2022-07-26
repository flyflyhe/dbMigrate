package controller

func init() {
	go func() {
		Start(":8080")
	}()
}
