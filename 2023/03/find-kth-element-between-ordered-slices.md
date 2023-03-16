# 寻找两个数组中第 K 大的数

## Question
给大两个从大到小排好序的数组 A 和 B，给定 K，找到 A 和 B 中第 K 大元素。

**func FindTopK(A, B []int, K int) (int, error)**

A=[10, 7, 3], B=[9, 6]

K=2, return 9

K=1, return 10

## Answer
```go
package main

import (
	"errors"
)

func FindTopK(A, B []int, K int) (int, error) {
	if K <= 0 {
		return 0, errors.New("K should be bigger than 0")
	}

	var topN int
	for ; K > 0; K-- {
		if len(A) <= 0 && len(B) <= 0 {
			return 0, errors.New("K is bigger than len(A)+len(B)")
		}

		if len(A) <= 0 || B[0] > A[0] {
			topN = B[0]
			B = B[1:]
		} else {
			topN = A[0]
			A = A[1:]
		}
	}
	return topN, nil
}
```

Unit test:
```go
package main

import "testing"

func TestFindTopK(t *testing.T) {
	cases := []struct {
		sliceA []int
		sliceB []int
		k int
		expectedNum int
		expectedError error
	} {
		{
			sliceA: []int{10,7,3},
			sliceB: []int{9,6},
			k: 2,
			expectedNum: 9,
			expectedError: nil,
		},
		{
			sliceA: []int{10,7,3},
			sliceB: []int{9,6},
			k: 1,
			expectedNum: 10,
			expectedError: nil,
		},
	}

	for _, c := range cases {
		ret, err := FindTopK(c.sliceA, c.sliceB, c.k)
		if ret != c.expectedNum || err != c.expectedError {
			t.Errorf("case: %+v, get num: %d, err: %v", c, ret, err)
		}
	}
}
```