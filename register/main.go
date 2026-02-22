/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package main

import (
	"fmt"
	"os"
	"strconv"
)

func trimQuotations(text string) string {
	size := len(text)
	if size > 1 {
		if text[0] == '"' && text[size-1] == '"' {
			return text[1 : size-1]
		}
		if text[0] == '\'' && text[size-1] == '\'' {
			return text[1 : size-1]
		}
	}
	return text
}

func getOptionString(args []string, key string) string {
	kLen := len(key)
	pos := len(args)
	var item string
	for pos > 0 {
		pos--
		item = args[pos]
		if item == key {
			return trimQuotations(args[pos+1])
		}
		if len(item) > kLen {
			// starts with "$(key)="
			if item[:kLen+1] == key+"=" {
				return trimQuotations(item[kLen+1:])
			}
		}
	}
	return ""
}
func getOptionInteger(args []string, key string) int {
	opt := getOptionString(args, key)
	num, err := strconv.Atoi(opt)
	if err == nil {
		return num
	} else {
		return 0
	}
}

func showHelp(path string) {
	fmt.Printf("\n"+
		"\n    Usages:"+
		"\n        %s <command> [options]"+
		"\n"+
		"\n    Commands:"+
		"\n        generate                Generate account."+
		"\n        modify                  Modify account info."+
		"\n        help                    Show help for commands."+
		"\n\n", path)
}

func doHelp(path string, args []string) {
	if len(args) == 1 {
		cmd := args[0]
		if cmd == "generate" {
			fmt.Printf("\n"+
				"\n    Usages:"+
				"\n        %s generate <type> [options]"+
				"\n"+
				"\n    Description:"+
				"\n        Generate account with type, e.g. 'USER', 'GROUP', 'STATION', 'ROBOT'."+
				"\n"+
				"\n    Generate Options:"+
				"\n        --seed <username>       Generate meta with seed string."+
				"\n        --founder <ID>          Generate group meta with founder ID."+
				"\n\n", path)
			return
		} else if cmd == "modify" {
			fmt.Printf("\n"+
				"\n    Usages:"+
				"\n        %s modify <ID> [options]"+
				"\n"+
				"\n    Descriptions:"+
				"\n        Modify account document with ID."+
				"\n"+
				"\n    Modify Options:"+
				"\n        --name <name>           Change name for user/group."+
				"\n        --avatar <URL>          Change avatar URL for user."+
				"\n        --host <IP>             Change IP for station."+
				"\n        --port <number>         Change port for station."+
				"\n        --owner <ID>            Change group info with owner ID."+
				"\n\n", path)
			return
		}
	}
	fmt.Printf("\n"+
		"\n    Usages:"+
		"\n        %s help <command>"+
		"\n"+
		"\n    Description:"+
		"\n        Show help for commands."+
		"\n"+
		"\n    Commands:"+
		"\n        generate"+
		"\n        modify"+
		"\n\n", path)
}

func main() {
	path := os.Args[0]
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "generate" {
			doGenerate(path, os.Args[2:])
			return
		} else if cmd == "modify" {
			doModify(path, os.Args[2:])
			return
		} else if cmd == "help" {
			doHelp(path, os.Args[2:])
			return
		}
	}
	showHelp(path)
}
