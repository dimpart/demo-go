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
package cpu

import (
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/sdk-go/cpu"
	. "github.com/dimpart/demo-go/sdk/client"
	. "github.com/dimpart/demo-go/sdk/common/db"
	. "github.com/dimpart/demo-go/sdk/common/protocol"
)

type LoginCommandProcessor struct {
	BaseCommandProcessor
}

// Override
func (cpu *LoginCommandProcessor) GetMessenger() IClientMessenger {
	messenger := cpu.BaseCommandProcessor.GetMessenger()
	cm, ok := messenger.(IClientMessenger)
	if ok {
		return cm
	}
	return nil
}

// private
func (cpu *LoginCommandProcessor) getDatabase() SessionDBI {
	messenger := cpu.GetMessenger()
	session := messenger.GetSession()
	return session.GetDatabase()
}

// Override
func (cpu *LoginCommandProcessor) ProcessContent(content Content, rMsg ReliableMessage) []Content {
	command, ok := content.(LoginCommand)
	if !ok {
		return nil
	}
	sender := command.ID()
	// save login command to session db
	db := cpu.getDatabase()
	if db.SaveLoginCommandMessage(sender, command, rMsg) {
		// OK
	} else {
		panic("failed to save login command")
	}
	// no need to response login command
	return nil
}
