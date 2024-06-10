package main

type Q struct {
	Query struct {
		Should struct {
			Match struct {
				Name string `json:"name"`
			} `json:"match"`
		} `json:"should"`
	} `json:"query"`
}

func main() {

	//r := gin.Default()
	//r.GET("/", func(c *gin.Context) {
	//	var zpr Q
	//	var s Q
	//	err := c.ShouldBindJSON(&zpr)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	data, _ := json.Marshal(zpr)
	//	_ = json.Unmarshal(data, &s)
	//	v := "qwe"
	//	log.Println(json.Marshal(v))
	//	log.Println(s.Query.Should.Match.Name)
	//	log.Println(zpr.Query.Should.Match.Name)
	//	c.Data(http.StatusOK, "application/json; charset=utf-8", data)
	//})
	//r.Run()
}
