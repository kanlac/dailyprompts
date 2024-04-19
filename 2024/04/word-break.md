# LeetCode 139. Word Break

https://leetcode.com/problems/word-break/description/

main_test.go
```go
package main

import (
	"testing"
)

func TestWordBreak(t *testing.T) {
	tests := []struct {
		s        string
		wordDict []string
		want     bool
	}{
		{"leetcode", []string{"leet", "code"}, true},
		{"applepenapple", []string{"apple", "pen"}, true},
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
		{"", []string{"any"}, true},
		{"abcd", []string{"ab", "abc", "cd"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			got := wordBreak(tt.s, tt.wordDict)
			if got != tt.want {
				t.Errorf("wordBreak(%q, %q) = %v, want %v", tt.s, tt.wordDict, got, tt.want)
			}
		})
	}
}
```

main.go
```go
package main

// dp[x] = true 表示字符串 s 中从第 0 个字符到第 x-1 个字符的部分可以被字典分词
func wordBreak(s string, wordDict []string) bool {
	var (
		l          = len(s)
		dp         = make([]bool, l+1)
		maxWordLen int
	)

	dp[0] = true // 认为空字符串可以被分词
		
	for _, w := range wordDict {
	  maxWordLen = max(maxWordLen, len(w))
	}

	// 在每次遍历中计算 s[:i] 是不是可以被字典分词
	// 每次 i++ 递增后都会基于前序结果计算
	for i := 1; i <= l; i++ {
		for j := i - 1; j >= max(i-maxWordLen-1, 0); j-- {
			// 如果 0:j 可以被字典分词，那么只需要看看 j:i 是不是字典中的某个词
			if dp[j] && contains(wordDict, s[j:i]) {
				dp[i] = true
				break
			}
		}
	}

	return dp[l]
}

func contains(words []string, target string) bool {
	for _, word := range words {
		if word == target {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

```
