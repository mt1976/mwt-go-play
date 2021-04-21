package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const constRateLen = 8

type fxRate struct {
	ccyPair   string
	bidRate   float64
	bidString string
	askRate   float64
	askString string
}

type fxRateCard struct {
	fxRates []fxRate
}

func main() {

	log.Println("INSTANCE")

	var rateCard []fxRate
	//eurjpy := getFXrate("EURJPY")
	//rate := getFXrate("EURJPY")
	rateCard = append(rateCard, getFXrate("AUDUSD"))
	rateCard = append(rateCard, getFXrate("EURGBP"))
	rateCard = append(rateCard, getFXrate("EURJPY"))
	rateCard = append(rateCard, getFXrate("EURUSD"))
	rateCard = append(rateCard, getFXrate("GBPUSD"))
	rateCard = append(rateCard, getFXrate("NZDUSD"))
	rateCard = append(rateCard, getFXrate("USDCAD"))
	rateCard = append(rateCard, getFXrate("USDCHF"))
	rateCard = append(rateCard, getFXrate("USDHKD"))
	rateCard = append(rateCard, getFXrate("USDJPY"))
	rateCard = append(rateCard, getFXrate("USDSGD"))
	//	rateCard = append(rateCard, getFXrate("XAUUSD"))

	//fmt.Println(rateCard, len(rateCard))
	log.Println("NEW RVMARKET GENERATED")
	f, err := os.Create("RVMARKET")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	outputString := ""
	for _, s := range rateCard {
		//fmt.Println(i, s)
		//fmt.Println(i, s)
		//fmt.Sprint(s)
		abc := fmt.Sprintf("s%ssptD%11sD%11s", s.ccyPair, s.bidString, s.askString)
		//fmt.Println(abc)
		outputString = outputString + abc + "\n"
	}
	log.Println("\n", outputString)
	_, err2 := f.WriteString(outputString)

	if err2 != nil {
		log.Fatal(err2)
	}

	log.Println("done")

	//	fmt.Println(q)
	//resp, err := ioutil.ReadAll(request.Body)
	//if err != nil {
	//	fmt.Println(err)
	//}

	//fmt.Println(resp)

	//	fc := forex.NewClient()
	//	params := map[string]string{"base": "EUR", "symbols": "GBP,USD"}
	//	rates, err := fc.Latest(params)
	//	if err != nil {
	//		log.Println(err)
	//	}

	//	fmt.Printf("rates = %+v", rates)
	// {Base:USD Date:2018-10-29 Rates:map[GBP:0.7801423425 EUR:0.8786574115]}

}

func findPos(inString string, searchString string) int {
	pos := strings.Index(inString, searchString)
	//fmt.Println("located", pos, len(inString), searchString, len(searchString))
	return pos
}

func findStartPos(inString string, searchString string) int {
	pos := findPos(inString, searchString)
	startPos := pos + len(searchString)
	return startPos
}

func getFXrate(inCCYpair string) fxRate {
	thisRate := fxRate{}
	thisRate.ccyPair = inCCYpair
	thisPair := "%5E" + inCCYpair
	url := fmt.Sprintf("https://www.barchart.com/forex/quotes/%s/overview", thisPair)
	//fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		log.Println(err, inCCYpair)
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err, inCCYpair)
		panic(err)
	}
	// show the HTML code as a string %s
	//fmt.Printf("%s\n", html)
	inString := string(html)

	searchString := "\"bidPrice\":\""
	searchString2 := "\",\"askPrice\":\""
	searchString3 := "\",\"bidSize\":\""

	bidPriceStart := findStartPos(inString, searchString)
	bidPriceStop := findPos(inString, searchString2)
	askPriceStart := findStartPos(inString, searchString2)
	askPriceStop := findPos(inString, searchString3)
	if askPriceStop == -1 {
		askPriceStop = askPriceStart + 7
	}

	thisRate.bidRate, _ = strconv.ParseFloat(inString[bidPriceStart:bidPriceStop], 64)
	thisRate.bidString = truncateString(inString[bidPriceStart:bidPriceStop], constRateLen)
	//fmt.Println("bidPrice=", bidPrice)

	thisRate.askRate, _ = strconv.ParseFloat(inString[askPriceStart:askPriceStop], 64)
	thisRate.askString = truncateString(inString[askPriceStart:askPriceStop], constRateLen)
	//fmt.Println("askPrice=", askPrice)
	log.Println(thisRate.ccyPair, bidPriceStart, askPriceStart, askPriceStop, thisRate.bidRate, thisRate.askRate)

	return thisRate
}

func truncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	// This code cannot support Japanese
	// orgLen := len(str)
	// if orgLen <= length {
	//     return str
	// }
	// return str[:length]

	// Support Japanese
	// Ref: Range loops https://blog.golang.org/strings
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}
