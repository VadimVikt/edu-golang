package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint
)

func TestCopySize(t *testing.T) {
	errDel := os.Remove("./testdata/tmp/output.txt")
	if errDel != nil {
		fmt.Println("Файл не найден", errDel)
	}
	input, err := os.Open("testdata/input.txt")
	if err != nil {
		t.Error(err)
	}
	defer input.Close()

	errCopy := Copy("testdata/input.txt", "testdata/tmp/output.txt", 0, 10_000)
	if errCopy != nil {
		t.Error(errCopy)
	}
	infoOut, errOut := os.Stat("testdata/tmp/output.txt")
	if errOut != nil {
		t.Error(errOut)
	}

	info, err := input.Stat()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(info.Size())

	require.Equal(t, info.Size(), infoOut.Size())
}

func TestLimit(t *testing.T) {
	t.Run("Offset 0, limit 0", func(t *testing.T) {
		var offset int64
		var limit int64
		err := Copy("testdata/input.txt", "testdata/output1.txt", offset, limit)
		if err != nil {
			t.Error(err)
		}
		expected, actual := getDataFromTests("./testdata/output1.txt", "testdata/out_offset0_limit0.txt")
		require.Equal(t, expected, actual)
		errDel := os.Remove("./testdata/output1.txt")
		if errDel != nil {
			fmt.Println("Файл не найден", errDel)
		}
	})

	t.Run("Offset 0, limit 10", func(t *testing.T) {
		var offset int64
		var limit int64 = 10
		err := Copy("testdata/input.txt", "testdata/output2.txt", offset, limit)
		if err != nil {
			t.Error(err)
		}
		expected, actual := getDataFromTests("./testdata/output2.txt", "testdata/out_offset0_limit10.txt")
		require.Equal(t, expected, actual)
		errDel := os.Remove("./testdata/output2.txt")
		if errDel != nil {
			fmt.Println("Файл не найден", errDel)
		}
	})

	t.Run("Offset 0, limit 10000", func(t *testing.T) {
		var offset int64
		var limit int64 = 10_000
		err := Copy("testdata/input.txt", "testdata/output3.txt", offset, limit)
		if err != nil {
			t.Error(err)
		}
		expected, actual := getDataFromTests("./testdata/output3.txt", "testdata/out_offset0_limit10000.txt")
		require.Equal(t, expected, actual)
		errDel := os.Remove("./testdata/output3.txt")
		if errDel != nil {
			fmt.Println("Файл не найден", errDel)
		}
	})
}

func TestOffset(t *testing.T) {
	t.Run("Offset 100, limit 1000", func(t *testing.T) {
		var offset int64 = 100
		var limit int64 = 1000
		err := Copy("testdata/input.txt", "testdata/output4.txt", offset, limit)
		if err != nil {
			t.Error(err)
		}
		expected, actual := getDataFromTests("./testdata/output4.txt", "testdata/out_offset100_limit1000.txt")
		require.Equal(t, expected, actual)
		errDel := os.Remove("./testdata/output4.txt")
		if errDel != nil {
			fmt.Println("Файл не найден", errDel)
		}
	})
	t.Run("Offset 6000, limit 1000", func(t *testing.T) {
		var offset int64 = 6000
		var limit int64 = 1000
		err := Copy("testdata/input.txt", "testdata/output5.txt", offset, limit)
		if err != nil {
			t.Error(err)
		}
		expected, actual := getDataFromTests("./testdata/output5.txt", "testdata/out_offset6000_limit1000.txt")
		require.Equal(t, expected, actual)
		errDel := os.Remove("./testdata/output5.txt")
		if errDel != nil {
			fmt.Println("Файл не найден", errDel)
		}
	})
}

func getDataFromTests(fileInput, fileOutput string) (in, out int64) {
	input, errIn := os.Open(fileInput)
	if errIn != nil {
		fmt.Printf("файл не найден %s ", errIn)
	}
	defer input.Close()
	output, errOut := os.Open(fileOutput)
	if errOut != nil {
		fmt.Printf("файл не найден %s ", errOut)
	}
	defer output.Close()
	infoIn, err := input.Stat()
	if err != nil {
		fmt.Printf("Нет информации о файле %s ", err)
	}
	infoOut, err := output.Stat()
	if err != nil {
		fmt.Printf("Нет информации о файле %s ", err)
	}
	return infoIn.Size(), infoOut.Size()
}
