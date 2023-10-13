package main

import (
	"context"
	"fmt"
	"homework/storage"
	"os"
	"path/filepath"
	"time"
)

func main() {

	td, err := os.MkdirTemp(os.TempDir(), "test-dir-*")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	defer os.Remove(td)

	d1 := filepath.Join(td, "dir1")
	err = os.Mkdir(d1, os.ModePerm)
	defer os.Remove(d1)

	f1 := filepath.Join(td, "test1.txt")
	err = os.WriteFile(f1, []byte("hello"), os.ModeTemporary)
	defer os.Remove(f1)

	f2 := filepath.Join(td, "test2.txt")
	err = os.WriteFile(f2, []byte("hello world"), os.ModeTemporary)
	defer os.Remove(f2)

	f3 := filepath.Join(d1, "test3.txt")
	err = os.WriteFile(f3, []byte("hello world from dir1"), os.ModeTemporary)
	defer os.Remove(f3)

	sizer := storage.NewSizer()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := sizer.Size(ctx, storage.NewLocalDir(td))
	fmt.Println(result)

	//c := make(chan any)
	//wg := sync.WaitGroup{}
	//m := sync.Mutex{}
	//
	//for i := 0; i < 3; i++ {
	//	wg.Add(1)
	//	go func(ch chan any, mut *sync.Mutex) {
	//		defer wg.Done()
	//
	//		select {
	//		case <-ch:
	//			mut.Lock()
	//			fmt.Println("ERROR OCCURRED")
	//			mut.Unlock()
	//		}
	//
	//	}(c, &m)
	//
	//}
	//c <- "DURA"
	//
	////wg.Add(1)
	////go func() {
	////	defer wg.Done()
	////
	////	c <- "DURA"
	////	return
	////}()
	//
	//wg.Wait()

}
