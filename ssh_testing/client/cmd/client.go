package main

import (
	"fmt"
	"github.com/VyacheslavIsWorkingNow/siv/ssh_testing/client/internal"
	"github.com/VyacheslavIsWorkingNow/siv/ssh_testing/client/internal/datagen"
	"log"
	"os"
	"strings"
	"time"
)

const (
	waterfall    = "waterfall"
	flashing     = "flashing"
	fifoFileName = "/Users/slavaruswarrior/Documents/GitHub/siv/visualizer/visfifo"
)

func main() {

	var userSchema string

	fmt.Printf("choose schema: waterfall | flashing\n")
	_, _ = fmt.Scanf("%s\n", &userSchema)

	if userSchema == waterfall {
		w := scanWaterfallSchema()
		errW := writeDataToFifo([]int{w.FlowSize, w.WaterfallSize, w.MaxAmplitude, w.Period})
		if errW != nil {
			return
		}
		generator := datagen.DataGenerator(
			time.Duration(w.Period)*time.Millisecond, w.FlowSize, w.MaxAmplitude, 0, 0, datagen.RandomMode)
		for i := 0; i < w.BatchCount; i++ {
			batch := <-generator
			err := writeDataToFifo(batch.Data)
			if err != nil {
				log.Fatalf("failed to send data to fifo in waterfall %e", err)
			}
		}
		errW = writeDataToFifo(generateEndData(w.FlowSize))
		if errW != nil {
			log.Fatalf("can't write end data")
		}
	} else if userSchema == flashing {
		f := scanFlashingSchema()
		errW := writeDataToFifo([]int{f.FlowSize, f.MaxAmplitude, f.Period})
		if errW != nil {
			log.Fatalf("can't write conf information")
		}
		generator := datagen.DataGenerator(
			time.Duration(f.Period)*time.Millisecond, f.FlowSize, f.MaxAmplitude, 0, 0, datagen.RandomMode)
		for i := 0; i < f.BatchCount; i++ {
			batch := <-generator
			err := writeDataToFifo(batch.Data)
			if err != nil {
				log.Fatalf("failed to send data to fifo in waterfall %e", err)
			}
		}
		errW = writeDataToFifo(generateEndData(f.FlowSize))
		if errW != nil {
			log.Fatalf("can't write end data")
		}
	} else {
		log.Fatalf("wrong graph type")
	}

	// canonicTest()

}

func scanWaterfallSchema() *internal.Waterfall {
	var flowSize, waterfallSize, maxAmplitude, period, batchCount int
	fmt.Println("Scan: flowSize | waterfallSize | maxAmplitude | period ms | batchCount")
	_, _ = fmt.Scanf("%d %d %d %d %d", &flowSize, &waterfallSize, &maxAmplitude, &period, &batchCount)
	w := &internal.Waterfall{
		FlowSize:      flowSize,
		WaterfallSize: waterfallSize,
		MaxAmplitude:  maxAmplitude,
		Period:        period,
		BatchCount:    batchCount,
	}

	return w
}

func scanFlashingSchema() *internal.Flashing {
	var flowSize, maxAmplitude, batchCount int
	var period float64
	fmt.Println("Scan: flowSize | maxAmplitude | period | batchCount")
	_, _ = fmt.Scanf("%d %d %d %f %d", &flowSize, &maxAmplitude, &period, &batchCount)
	return nil
}

func writeDataToFifo(data []int) error {
	resultString := strings.Join(strings.Fields(fmt.Sprint(data)), " ")
	resultString = resultString[1 : len(resultString)-1]

	fmt.Println("send batch:", resultString)
	fifoFile, err := os.OpenFile(fifoFileName, os.O_WRONLY|os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		return fmt.Errorf("failed to open fifo file %w", err)
	}
	defer func(fifoFile *os.File) {
		_ = fifoFile.Close()
	}(fifoFile)

	_, err = fifoFile.WriteString(resultString)
	if err != nil {
		return fmt.Errorf("failed to write fifo file %w", err)
	}
	return nil
}

func generateEndData(N int) []int {
	arr := make([]int, N)
	arr[0] = -1
	return arr
}

func canonicTest() {
	w := &internal.Waterfall{
		FlowSize:      50,
		WaterfallSize: 20,
		MaxAmplitude:  70,
		Period:        250,
		BatchCount:    30,
	}
	errW := writeDataToFifo([]int{w.FlowSize, w.WaterfallSize, w.MaxAmplitude, w.Period})
	if errW != nil {
		return
	}
	for i := 0; i < w.BatchCount; i++ {
		err := writeDataToFifo(datagen.CanonicData[i])
		if err != nil {
			log.Fatalf("failed to send data to fifo in waterfall %e", err)
		}
		time.Sleep(time.Millisecond * time.Duration(w.Period))
	}
	errW = writeDataToFifo(generateEndData(w.FlowSize))
	if errW != nil {
		log.Fatalf("can't write end data")
	}
}
