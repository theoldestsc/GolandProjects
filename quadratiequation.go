package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

const pageTop1  = `<h1>Quadratic Equation Solver</h1>`
const anError1 = `<p class="error">%s</p>`
const form1 = `
<h><b>Quadratic Equation Solver</b></h1>
<form action="/" method="POST">
<label for="numbers">Solves equations of the form ax<sup>2</sup> + bx + c</label><br />
<input type="text" name="a" size="5">
<a>x <sup>2</sup> +</a>
<input type="text" name="b" size="5">
<a>x +</a>
<input type="text" name="c" size="5">
<a>&rarr;</a>
<input type="submit" value="Calculate">
</form>`

//Расчеты
//TODO:Анализ дискриминанта
//TODO:Возврат значния
func QuadraticEquation(a,b,c float64) []string{
	var results []string

	/*sA:=strconv.FormatFloat(a,'E',0,64)
	mas := strings.Split(sA,"E")
	fmt.Println(mas)
	fmt.Println(sA)*/
	//formatResult:=""

	fmt.Printf("Your Quadratic:%fx^2+%fx+%f=0\n",a,b,c)
	Discriminant :=math.Pow(b,2) - 4*a*c
	fmt.Println(Discriminant)
	if Discriminant > 0{
		x1:=(-b+math.Sqrt(Discriminant))/(2*a)
		x2:=(-b-math.Sqrt(Discriminant))/(2*a)
		str1:=""+fmt.Sprint(x1)
		str2:=""+fmt.Sprint(x2)

		results = append(results,str1,str2)
	}
	if Discriminant == 0{
		x:=-b/2*a
		str:=""+fmt.Sprint(x)

		results = append(results,str)
	}
	if Discriminant < 0{
		x1:=(-b+math.Sqrt(-Discriminant))/(2*a)
		x2:=(-b-math.Sqrt(-Discriminant))/(2*a)
		str1:=""+fmt.Sprint(x1)+"i"
		str2:=""+fmt.Sprint(x2)+"i"

		results = append(results,str1,str2)
	}

	return results
}

/*
func main()  {
	var a float64
	var b float64
	var c float64
	fmt.Printf("Введите {number}x^2: Введите {number}x: Введите c: ")
	fmt.Scanf("%f %f %f",&a, &b, &c)
	fmt.Println(a,b,c)
	fmt.Println(QuadraticEquation(a,b,c))
}
*/



//Сервер
func main(){
	//URL удет "/" - выполнить эту функцию
	http.HandleFunc("/",homePage1)
	//СЛУШАЕМ ПОРТ
	if err := http.ListenAndServe(":9001", nil); err != nil{
		log.Fatal("failed to start server", err)
	}
}

func processRequest1(request *http.Request) (map[string]float64,string,bool) {
	abc := make(map[string]float64)
	flag_a:=false
	flag_b:=false
	flag_c:=false
	if slice1,found1 := request.Form["a"];found1 && len(slice1)>0 {
		text_a := slice1[0]
		flag_a = found1
		fmt.Println(text_a, flag_a)
		if value_a, err := strconv.ParseFloat(text_a, 64); err != nil {
			return abc, "'" + text_a + "' is invalid", false
		}else{
			abc["a"]=value_a
		}

	}
	if slice1,found1 := request.Form["b"];found1 && len(slice1)>0{
		text_b := slice1[0]
		flag_b = found1
		fmt.Println(text_b, flag_b)
		if value_b, err := strconv.ParseFloat(text_b, 64); err != nil {
			return abc, "'" + text_b + "' is invalid", false
		}else{
			abc["b"]=value_b
		}
	}
	if slice1,found1 := request.Form["c"];found1 && len(slice1)>0{
		text_c := slice1[0]
		flag_c = found1
		fmt.Println(text_c, flag_c)
		if value_c, err := strconv.ParseFloat(text_c, 64); err != nil {
			return abc, "'" + text_c + "' is invalid", false
		}else{
			abc["c"]=value_c
		}
	}
	if len(abc) == 0{
		return abc,"",false // при первом отображении данные отсутствуют
	}
	return abc, " ",true
}
func formatEquation(a,b,c float64, x1,x2 string)string{
	return fmt.Sprintf(`<h>%fx<sup>2</sup> + %fx + %f &rarr;x = (%s) or x = (%s)</h>`,a,b,c,x1,x2)
}
func homePage1(writer http.ResponseWriter,request *http.Request) {
	//http.ServeFile(writer, request, "static/main.html") - отобразить html файл
	fmt.Fprintf(writer,pageTop1,form1)
	err := request.ParseForm()
	if err != nil {
		fmt.Fprintf(writer, anError1, err)
	} else {
		if dictionary, message, ok := processRequest1(request); ok {
			equation := QuadraticEquation(dictionary["a"],dictionary["b"],dictionary["c"])
			fmt.Println(equation,message)
			fmt.Fprintf(writer, formatEquation(dictionary["a"],dictionary["b"],dictionary["c"],equation[0],equation[1]))
		}else if message !=""{
			fmt.Fprintf(writer,anError1,message)
			}
		}
	}