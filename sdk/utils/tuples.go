/* license: https://mit-license.org
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2026 Albert Moky
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
package utils

type Pair[A, B any] struct {
	First  A
	Second B
}

type Triplet[A, B, C any] struct {
	First  A
	Second B
	Third  C
}

//type Quartet[A, B, C, D any] struct {
//	First  A
//	Second B
//	Third  C
//	Fourth D
//}
//
//type Quintet[A, B, C, D, E any] struct {
//	First  A
//	Second B
//	Third  C
//	Fourth D
//	Fifth  E
//}

//
//  Creation
//

func NewPair[A, B any](first A, second B) *Pair[A, B] {
	return &Pair[A, B]{
		First:  first,
		Second: second,
	}
}

func NewTriplet[A, B, C any](first A, second B, third C) *Triplet[A, B, C] {
	return &Triplet[A, B, C]{
		First:  first,
		Second: second,
		Third:  third,
	}
}

//func NewQuartet[A, B, C, D any](first A, second B, third C, fourth D) *Quartet[A, B, C, D] {
//	return &Quartet[A, B, C, D]{
//		First:  first,
//		Second: second,
//		Third:  third,
//		Fourth: fourth,
//	}
//}
//
//func NewQuintet[A, B, C, D, E any](first A, second B, third C, fourth D, fifth E) *Quintet[A, B, C, D, E] {
//	return &Quintet[A, B, C, D, E]{
//		First:  first,
//		Second: second,
//		Third:  third,
//		Fourth: fourth,
//		Fifth:  fifth,
//	}
//}
