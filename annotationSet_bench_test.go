package lingo

import (
	"sort"
	"testing"
)

func (as AnnotationSet) index2(a *Annotation) int {
	sort.Sort(as)
	f := func(i int) bool { return as[i] == a }
	return sort.Search(len(as), f)
}

var benchIndexRes int

func benchASIndex(size int, b *testing.B) {
	var as AnnotationSet
	for i := 0; i < size; i++ {
		as = append(as, new(Annotation))
	}

	doesntcontain := new(Annotation)
	contains := as[0]

	for n := 0; n < b.N; n++ {
		benchIndexRes = as.Index(doesntcontain)
		benchIndexRes = as.Index(contains)
	}
}

func benchASIndex2(size int, b *testing.B) {
	var as AnnotationSet
	for i := 0; i < size; i++ {
		as = append(as, new(Annotation))
	}

	doesntcontain := new(Annotation)
	contains := as[0]

	for n := 0; n < b.N; n++ {
		benchIndexRes = as.index2(doesntcontain)
		benchIndexRes = as.index2(contains)
	}
}

func BenchmarkAnnotationSetIndex_1(b *testing.B)    { benchASIndex(1, b) }
func BenchmarkAnnotationSetIndex_2(b *testing.B)    { benchASIndex(2, b) }
func BenchmarkAnnotationSetIndex_8(b *testing.B)    { benchASIndex(8, b) }
func BenchmarkAnnotationSetIndex_16(b *testing.B)   { benchASIndex(16, b) }
func BenchmarkAnnotationSetIndex_32(b *testing.B)   { benchASIndex(32, b) }
func BenchmarkAnnotationSetIndex_64(b *testing.B)   { benchASIndex(64, b) }
func BenchmarkAnnotationSetIndex_128(b *testing.B)  { benchASIndex(128, b) }
func BenchmarkAnnotationSetIndex_256(b *testing.B)  { benchASIndex(256, b) }
func BenchmarkAnnotationSetIndex_512(b *testing.B)  { benchASIndex(512, b) }
func BenchmarkAnnotationSetIndex_1024(b *testing.B) { benchASIndex(1024, b) }

func BenchmarkAnnotationSetIndex2_1(b *testing.B)    { benchASIndex2(1, b) }
func BenchmarkAnnotationSetIndex2_2(b *testing.B)    { benchASIndex2(2, b) }
func BenchmarkAnnotationSetIndex2_8(b *testing.B)    { benchASIndex2(8, b) }
func BenchmarkAnnotationSetIndex2_16(b *testing.B)   { benchASIndex2(16, b) }
func BenchmarkAnnotationSetIndex2_32(b *testing.B)   { benchASIndex2(32, b) }
func BenchmarkAnnotationSetIndex2_64(b *testing.B)   { benchASIndex2(64, b) }
func BenchmarkAnnotationSetIndex2_128(b *testing.B)  { benchASIndex2(128, b) }
func BenchmarkAnnotationSetIndex2_256(b *testing.B)  { benchASIndex2(256, b) }
func BenchmarkAnnotationSetIndex2_512(b *testing.B)  { benchASIndex2(512, b) }
func BenchmarkAnnotationSetIndex2_1024(b *testing.B) { benchASIndex2(1024, b) }
