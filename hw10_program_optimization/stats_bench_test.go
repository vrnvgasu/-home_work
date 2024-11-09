package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s -cpuprofile=cpu_old.out -memprofile=mem_old.out .
// go tool pprof -http=":8090" cpu_old.out mem_old.out
// go test -bench=BenchmarkGetDomainStat -benchmem -benchtime 10s -cpuprofile=cpu_new.out -memprofile=mem_new.out .
// go tool pprof -http=":8090" cpu_new.out mem_new.out

func BenchmarkGetDomainStat(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) { // запуск параллельно
		for pb.Next() { // крутим внутри параллельного запуска цикл for
			r, err := zip.OpenReader("testdata/users.dat.zip")
			if err != nil {
				b.Fatal(err)
			}
			defer r.Close()

			data, err := r.File[0].Open()
			if err != nil {
				b.Fatal(err)
			}
			_, err = GetDomainStat(data, "biz")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkGetDomainStatOld(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) { // запуск параллельно
		for pb.Next() { // крутим внутри параллельного запуска цикл for
			r, err := zip.OpenReader("testdata/users.dat.zip")
			if err != nil {
				b.Fatal(err)
			}
			defer r.Close()

			data, err := r.File[0].Open()
			if err != nil {
				b.Fatal(err)
			}
			_, err = GetDomainStatOld(data, "biz")
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
