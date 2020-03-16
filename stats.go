package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const form string  = `<form action="/" method="POST">
<label for="numbers">Numbers (comma or space-separated):</label><br />"
<input type="text" name="numbers" size="60"><br />"
<input type="submit" value="Calculate">
</form>`
const anError = `<p class="error">%s</p>`
const pageTop  = `<h1>Title</h1>`
const pageBottom  = `<h2>Under Title</h2>`
//Простые вычисления
type statistics struct{
	numbers []float64
	mean float64
	median float64
	std_dev float64
	modal float64
	module_err error
}

func getStats(numbers []float64) (stats statistics){
	stats.numbers = numbers
	sort.Float64s(stats.numbers)
	stats.mean = sum(numbers)/float64(len(numbers))
	stats.median = median(numbers)
	stats.std_dev = stdDev(numbers,stats.mean)
	stats.modal,stats.module_err = modal(numbers)
	return stats
}

func stdDev(numbers []float64, mean float64)float64{
	var values []float64
	for _,value:=range numbers{
		x :=math.Pow(value-mean, 2)
		values = append(values,x)
	}
	v := math.Sqrt(sum(values)/float64(len(numbers)-1))
	return v
}

func median(numbers []float64)float64{
	middle := len(numbers)/2
	result := numbers[middle]
	if len(numbers)%2 == 0{
		result = (result + numbers[middle-1])/2
	}
	return result
}

func sum(numbers []float64) (total float64) {
	for _,x := range numbers{
		total += x
	}
	return total
}

func modal(numbers []float64) (number float64,err error) {
	repeats := make(map[float64]int)
	err = errors.New("There is no module")
	for _, number1 := range numbers {
		if _, ok := repeats[number1]; !ok {
			repeats[float64(number1)] = 1
		} else {
			repeats[float64(number1)] += 1
		}
	}
	value := 1
	entered := 0
	var numberRepeat float64 = -1
	for number := range repeats {
		if repeats[number] > value {
			value = repeats[number]
			numberRepeat = number
		}
		if repeats[number] == value {
			entered+=1
			continue
		}

	}
	return numberRepeat, nil
}
//Анализ формы
func processRequest(request *http.Request) ([]float64,string,bool) {
	var numbers []float64
	if slice, found := request.Form["numbers"]; found && len(slice) > 0 {
		text := strings.Replace(slice[0], ",", " ", -1)
		for _, field := range strings.Fields(text) {
			if x, err := strconv.ParseFloat(field, 64); err != nil {
				return numbers, "'" + field + "' is invalid", false
			} else {
				numbers = append(numbers, x)
			}
		}
	}
	if len(numbers) == 0{
		return numbers,"",false // при первом отображении данные отсутствуют
	}
	return numbers,"",true
}

func formatStats(stats statistics)string{
	return fmt.Sprintf(`<table border = "1">
			<tr><th colspan = "2">Results</th></tr>
			<tr><td>Numbers</td><td>%v</td></tr>
			<tr><td>Count</td><td>%d</td></tr>
			<tr><td>Mean</td><td>%f</td></tr>
			<tr><td>Median</td><td>%v</td></tr>
			<tr><td>Std.Dev</td><td>%v</td></tr>
			<tr><td>Modal</td><td>%v</td></tr>`, stats.numbers,len(stats.numbers),stats.mean,stats.median,stats.std_dev,stats.modal)
}




func main(){
	//URL удет "/" - выполнить эту функцию
	http.HandleFunc("/",homePage)
	//СЛУШАЕМ ПОРТ
	if err := http.ListenAndServe(":9001", nil); err != nil{
		log.Fatal("failed to start server", err)
	}

}

func homePage(writer http.ResponseWriter,request *http.Request){
	err:=request.ParseForm()
	//fmt.Println(writer)
	fmt.Fprintf(writer,pageTop,form)
	//fmt.Println(request)
	//fmt.Println(*request)

	//fmt.Fprintln(writer,request.FormValue("numbers")) Распарсил форму, вывел значение формы

	if err != nil{
		fmt.Fprintf(writer, anError,err)
	}else{
		if numbers,message,ok := processRequest(request);ok{
			stats := getStats(numbers)
			fmt.Fprintf(writer, formatStats(stats))
		}else if message !=""{
			fmt.Fprintf(writer,anError,message)
		}
	}
	fmt.Fprintf(writer,pageBottom)
}