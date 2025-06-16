package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"github.com/wcharczuk/go-chart/v2"
	"sort"
)

type JobLogs struct{
	Jobname string
	TimeStamp time.Time
	Status string
}

var (
	startmap = make(map[string]time.Time)
	endmap   = make(map[string]time.Time)
	bothstartandend = make(map[string]time.Duration)
	onlystarttime = make(map[string]time.Time)
	onlyendtime = make(map[string]time.Time)
	invalidjob =make(map[string]time.Duration)
)

func Piechart(val1,val2,val3,val4 int){
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: []chart.Value{
			{Value: float64(val1), Label: "both start and end time"},
			{Value: float64(val2), Label: "only start time"},
			{Value: float64(val3), Label: "only end time"},
			{Value: float64(val4), Label: "invalid duration"},
		},
	}
	f, err := os.Create("output_pie.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	pie.Render(chart.PNG, f)
}



// func BarChartForDurations(data map[string]time.Duration) {
// 	// Extract and sort keys (job names) for consistent order
// 	keys := make([]string, 0, len(data))
// 	for k := range data {
// 		keys = append(keys, k)
// 	}
// 	sort.Strings(keys)

// 	// Prepare bars: convert durations to seconds (float64)
// 	bars := make([]chart.Value, len(keys))
// 	for i, k := range keys {
// 		bars[i] = chart.Value{
// 			Value: data[k].Seconds(),
// 			Label: k,
// 		}
// 	}

// 	barChart := chart.BarChart{
// 		Width:  1024,
// 		Height: 512,
// 		Bars:   bars,
// 		YAxis: chart.YAxis{
// 			Name: "Duration (seconds)",
// 			// Minimal styling:
// 			// Style: chart.Style{
// 			// 	Show: true,
// 			// },
// 		},
// 		// XAxis: chart.Style{
// 		// 	Show: true,
// 		// },
// 	}

// 	f, err := os.Create("bothstartandend_barchart.png")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()

// 	err = barChart.Render(chart.PNG, f)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func BarChartForDurations(data map[string]time.Duration) {
	if len(data) == 0 {
		fmt.Println("No data provided for bar chart; skipping chart generation")
		return
	}

	// Extract and sort keys
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Prepare bars and check if all zero
	allZero := true
	bars := make([]chart.Value, len(keys))
	for i, k := range keys {
		sec := data[k].Seconds()
		bars[i] = chart.Value{
			Value: sec,
			Label: k,
		}
		if sec != 0 {
			allZero = false
		}
	}

	if allZero {
		fmt.Println("All duration values are zero; skipping bar chart generation")
		return
	}

	barChart := chart.BarChart{
		Width:  1024,
		Height: 512,
		Bars:   bars,
		YAxis: chart.YAxis{
			Name: "Duration (seconds)",
			Style: chart.Style{
				FontSize: 10,
			},
		},
		// XAxis: chart.XAxis{
		// 	Style: chart.Style{
		// 		FontSize: 10,
		// 	},
		// },
	}

	f, err := os.Create("bothstartandend_barchart.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = barChart.Render(chart.PNG, f)
	if err != nil {
		panic(err)
	}
}


func main() {
	fmt.Println("graph visualization...")
	// file,err:=os.Open("logfile.txt")
	file,err:=os.Open("logfile2.txt")
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer file.Close()
	scanner :=bufio.NewScanner(file)

	var enteries []JobLogs

	for scanner.Scan(){
		line:=scanner.Text()
		fields:=strings.Split(line,",")
		if len(fields)!=3{
			fmt.Println("invalid log format",line)
			continue
		}
		jobname:=strings.TrimSpace(fields[0])
		timestampstr:=strings.TrimSpace(fields[1])
		status:=strings.TrimSpace(fields[2])
		timestamp,err:=time.Parse(time.RFC3339,timestampstr)

		if err!=nil{
			fmt.Println("cannot convert timestamp string to time.Time",timestampstr)
			continue
		}

		if status=="start"{
			startmap[jobname]=timestamp
		}else if status=="end"{
			endmap[jobname]=timestamp
		}else{
			fmt.Println("invalid status",status)
			continue
		}

		entry:=JobLogs{
			Jobname: jobname,
			TimeStamp: timestamp,
			Status: status,
		}
		enteries=append(enteries, entry)
	}
	if err:=scanner.Err();err!=nil{
		fmt.Println(err)
		return
	}
	// for _,entry:=range enteries{
	// 	fmt.Println(entry.Jobname,entry.TimeStamp,entry.Status)
	// }
	// fmt.Println(enteries[0].TimeStamp.Sub(enteries[1].TimeStamp))

	// fmt.Println("start map")
	// for jbname,start:=range startmap{
	// 	fmt.Println(jbname,start)
	// }
	// fmt.Println("end map")
	// for jbname,end:=range endmap{
	// 	fmt.Println(jbname,end)
	// }

	for _,y:= range enteries{
		z:=y.Jobname
		// fmt.Print(z,"->")
		_,exists1:=startmap[z] 
		_,exists2:=endmap[z]
		if exists1 && exists2{
			// fmt.Println("start time",startmap[z],"end time",endmap[z])
			duration := endmap[z].Sub(startmap[z])
			if duration<0{
				invalidjob[z]=duration
			}else{
				// fmt.Println(z,"duration","->",duration)
				bothstartandend[z]=duration
			}
		}else if exists1{
			// fmt.Println("only start time present for",z,"->",startmap[z])
			onlystarttime[z]=startmap[z]
		}else{
			// fmt.Println("only end time present",z,"->",endmap[z])
			onlyendtime[z]=endmap[z]
		}
	}


	if len(bothstartandend)>0{
		fmt.Println("BOTH START AND END TIME")
		fmt.Println("count of valid data is",len(bothstartandend))
		for key,val:= range bothstartandend{
			fmt.Println(key,val)
		}
	}
		if len(onlystarttime)>0{
			fmt.Println("START TIME ONLY")
			fmt.Println("count of only start time data is",len(onlystarttime))
			for key,val:= range onlystarttime{
				fmt.Println(key,val)
			}
		}
		if len(onlyendtime)>0{	
			fmt.Println("END TIME ONLY")
			fmt.Println("count of only end time data is",len(onlyendtime))
			for key,val:= range onlyendtime{
			fmt.Println(key,val)
		}
		}
		if len(invalidjob)>0{	
			fmt.Println("INVALID JOBS")
			fmt.Println("count of invalid data is",len(invalidjob))
			for key,val:= range invalidjob{
			fmt.Println(key,val)
		}
		}
		Piechart(len(bothstartandend),len(onlystarttime),len(onlyendtime),len(invalidjob))
		BarChartForDurations(bothstartandend)

}