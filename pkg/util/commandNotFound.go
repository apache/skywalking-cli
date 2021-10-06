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

	"github.com/urfave/cli/v2"
)

const GAP = 3 // Control the appropriate edit distance.

// CommandNotFound is executed when the command entered does not exist.
func CommandNotFound(c *cli.Context, s string) {
	suppose := make([]string, 0)
	var parentCommand string
	if len(c.Lineage()) == 1 {
		parentCommand = "swctl"
	} else {
		parentCommand = c.Lineage()[1].Args().First()
	}
	fmt.Printf("Error: unknown command \"%s\" for \"%s\" \n\n", s, parentCommand)
	// Record commands whose edit distance is less than GAP to suppose.
	for index := range c.App.Commands {
		commandName := c.App.Commands[index].Name
		distance := minEditDistance(commandName, s)
		if distance <= GAP && commandName != "help" {
			suppose = append(suppose, commandName)
		}
	}
	if len(suppose) != 0 {
		fmt.Println("Do you mean this?")
		for index := range suppose {
			fmt.Printf("\t%s\n", suppose[index])
		}
		fmt.Println()
	}
	if len(c.Lineage()) == 1 {
		fmt.Printf("Run '%s --help' for usage.\n", parentCommand)
	} else {
		fmt.Printf("Run 'swctl %s --help' for usage.\n", parentCommand)
	}
}

// minEditDistance calculates the edit distance of two strings.
func minEditDistance(word1, word2 string) int {
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
				dp[i][j] = min(dp[i][j-1], dp[i-1][j], dp[i-1][j-1]) + 1
			}
		}
	}
	return dp[m][n]
}

// min get The minimum of the args.
func min(args ...int) int {
	min := args[0]
	for _, item := range args {
		if item < min {
			min = item
		}
	}
	return min
}
