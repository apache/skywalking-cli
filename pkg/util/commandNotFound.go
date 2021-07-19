// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package util

import (
	"fmt"

	"github.com/urfave/cli"
)

const GAP = 3

// Handle error through CommandNotFound function when the input command does not exist
func CommandNotFound(c *cli.Context, s string) {
	suppose := make([]string, 0)

	fmt.Printf("Command '%s' not found. ", s)
	for index := range c.App.Commands {
		commandName := c.App.Commands[index].Name
		distance := minDistance(commandName, s)
		// Record commands that edit distance is less than 3
		if distance <= GAP {
			suppose = append(suppose, commandName)
		}
	}

	if len(suppose) != 0 {
		fmt.Println("Do you mean:")
		for index := range suppose {
			if c.Parent() == nil {
				fmt.Printf("\t%s\n", suppose[index])
			} else {
				fmt.Printf("\t%s %s\n", c.Parent().Args()[0], suppose[index])
			}
		}
	}
	fmt.Println()
}

// Calculate the edit distance of two strings
func minDistance(word1, word2 string) int {
	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 0; i < m+1; i++ {
		dp[i][0] = i
	}
	for j := 0; j < n+1; j++ {
		dp[0][j] = j
	}
	for i := 1; i < m+1; i++ {
		for j := 1; j < n+1; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = Min(dp[i][j-1], dp[i-1][j], dp[i-1][j-1]) + 1
			}
		}
	}
	return dp[m][n]
}
func Min(args ...int) int {
	min := args[0]
	for _, item := range args {
		if item < min {
			min = item
		}
	}
	return min
}
