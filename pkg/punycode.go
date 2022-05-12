package pkg

import (
	"math"
	"regexp"
	"strings"
	"unicode/utf8"
)

type Punycode struct{}

// Highest positive signed 32-bit float value
const maxInt = int32(2147483647) // aka. 0x7FFFFFFF or 2^31-1

const base = int32(36)
const tMin = int32(1)
const tMax = int32(26)
const skew = int32(38)
const damp = int32(700)
const initialBias = int32(72)
const initialN = int32(128) // 0x80
const delimiter = "-"       // '\x2D'

/** Convenience shortcuts */
const baseMinusTMin = base - tMin

var regexPunycodeRe = regexp.MustCompile("^xn--")
var regexNonASCIIRe = regexp.MustCompile("[^\\0-\x7E]")
var regexSeparators = regexp.MustCompile("[\x2E\u3002\uFF0E\uFF61]")

func Map(array []string, fn func(string) string) []string {

	var n = len(array)
	var result = make([]string, n)

	for n > 0 {
		n--
		result[n] = fn(array[n])
	}
	return result
}

func MapDomain(str string, fn func(string) string) string {
	var parts = strings.Split(str, "@")
	var result = ""

	if len(parts) > 1 {
		// In email addresses, only the domain name should be punycoded. Leave
		// the local part (i.e. everything up to `@`) intact.
		result = parts[0] + "@"
		str = parts[1]
	}

	// Avoid `split(regex)` for IE8 compatibility. See #17.
	str = regexSeparators.ReplaceAllString(str, "\x2E")
	var labels = strings.Split(str, ".")
	var encoded = strings.Join(Map(labels, fn)[:], ".")
	return result + encoded
}

func (p *Punycode) ToASCII(input string) string {
	return MapDomain(input, func(s string) string {
		if regexNonASCIIRe.MatchString(s) {
			// TODO: Handle error
			return "xn--" + p.Encode(s)
		}
		return s
	})
}

func (p *Punycode) ToUnicode(input string) string {
	return MapDomain(input, func(s string) string {
		if regexPunycodeRe.MatchString(s) {
			// TODO: Handle error
			slice := strings.ToLower(Slice(s, 4, utf8.RuneCountInString(s)))
			return p.Decode(slice)
		}
		return s
	})
}

func (p *Punycode) Ucs2Decode(input string) string {
	// TODO
	return input
}

func (p *Punycode) Ucs2Encode(array []rune) string {
	return string(array)
}

func (p *Punycode) Adapt(delta int32, numPoints int32, firstTime bool) int32 {
	var k int32 = 0

	delta = int32(math.Floor(float64(delta) / float64(damp)))
	if !firstTime {
		delta = delta >> 1
	}

	delta += int32(float64(delta) / float64(numPoints))

	for ; /* no initialization */ delta > baseMinusTMin*tMax>>1; k += base {
		delta = int32(math.Floor(float64(delta) / float64(baseMinusTMin)))
	}

	return int32(math.Floor(float64(k+(baseMinusTMin+1)*delta) / float64(delta+skew)))
}

func (p *Punycode) BasicToDigit(codePoint int32) int32 {
	if codePoint-0x30 < 0x0A {
		return codePoint - 0x16
	}
	if codePoint-0x41 < 0x1A {
		return codePoint - 0x41
	}
	if codePoint-0x61 < 0x1A {
		return codePoint - 0x61
	}
	return base
}

func (p *Punycode) DigitToBasic(digit int32, flag int) int32 {

	var isFlagEnabled int32 = 1
	var isDigitSmaller int32 = 1

	if digit >= 26 {
		isDigitSmaller = 0
	}

	if flag == 0 {
		isFlagEnabled = 0
	}

	return digit + 22 + 75*isDigitSmaller - ((isFlagEnabled) << 5)
}

func (p *Punycode) Encode(input string) string {
	// TODO
	var output []string

	input = p.Ucs2Decode(input)
	var inputLength = utf8.RuneCountInString(input)

	// Initialize the state.
	var n = int32(initialN)
	var delta = int32(0)
	var bias = int32(initialBias)

	// Handle the basic code points.
	for _, currentValue := range input {
		if currentValue < 0x80 {
			output = append(output, string(currentValue))
		}
	}

	var basicLength = len(output)
	var handledCPCount = basicLength

	// `handledCPCount` is the number of code points that have been handled;
	// `basicLength` is the number of basic code points.

	// Finish the basic string with a delimiter unless it's empty.
	if basicLength > 0 {
		output = append(output, delimiter)
	}

	// Main encoding loop:
	for handledCPCount < inputLength {
		// All non-basic code points < n have been handled already. Find the next
		// larger one:
		var m = int32(maxInt)
		for _, currentValue := range input {
			if currentValue >= n && currentValue < m {
				m = currentValue
			}
		}

		// Increase `delta` enough to advance the decoder's <n,i> state to <m,0>,
		// but guard against overflow.
		var handledCPCountPlusOne = int32(handledCPCount + 1)
		var change = float64((maxInt - delta) / handledCPCountPlusOne)
		if float64(m-n) > math.Floor(change) {
			// TODO: Return error here
			//return "", errors.New("overflow: input needs wider integers to process")
		}

		delta += (m - n) * int32(handledCPCountPlusOne)
		n = m

		for _, currentValue := range input {
			delta++
			if currentValue < n && delta > maxInt {
				// TODO: Return error here
				// return "", errors.New("overflow: input needs wider integers to process")
			}

			if currentValue == n {
				var q = delta

				for k := base; ; /* no condition */ k += base {
					var t int32

					if int32(k) <= bias {
						t = tMin
					} else {
						if int32(k) >= bias+tMax {
							t = tMax
						} else {
							t = int32(k) - bias
						}
					}

					if q < t {
						break
					}

					var qMinusT = q - t
					var baseMinusT = base - t

					output = append(output, string(rune(p.DigitToBasic(t+qMinusT%baseMinusT, 0))))
					q = int32(math.Floor(float64(qMinusT) / float64(baseMinusT)))
				}
				output = append(output, string(rune(p.DigitToBasic(q, 0))))
				bias = p.Adapt(delta, handledCPCountPlusOne, handledCPCount == basicLength)
				delta = 0
				handledCPCount++
			}
		}
		delta++
		n++
	}

	return strings.Join(output[:], "")
}

func (p *Punycode) Decode(input string) string {
	// Don't use UCS-2.
	var output []int32
	var inputLength = utf8.RuneCountInString(input)
	var i = int32(0)
	var n = int32(initialN)
	var bias = int32(initialBias)

	// Handle the basic code points: let `basic` be the number of input code
	// points before the last delimiter, or `0` if there is none, then copy
	// the first basic code points to the output.

	var basic = strings.LastIndex(input, delimiter)
	if basic < 0 {
		basic = 0
	}

	for j := 0; j < basic; {
		j++
		// if it's not a basic code point
		if CharCodeAt(input, j) >= 0x80 {
			// TODO: Handle errors
			//return "", errors.New("overflow: input needs wider integers to process")
		}
		output = append(output, CharCodeAt(input, j))
	}

	// Main decoding loop: start just after the last delimiter if any basic code
	// points were copied; start at the beginning otherwise.
	var index = basic + 1

	if basic <= 0 {
		index = 0
	}

	for index < inputLength {
		// `index` is the index of the next character to be consumed.
		// Decode a generalized variable-length integer into `delta`,
		// which gets added to `i`. The overflow checking is easier
		// if we increase `i` as we go, then subtract off its starting
		// value at the end to obtain `delta`.

		var oldi = i
		var k = base

		for w := 1; ; /* no condition */ k += base {
			if index >= inputLength {
				//return "", errors.New("overflow: input needs wider integers to process")

			}

			var digit = p.BasicToDigit(CharCodeAt(input, index))
			index++

			if digit >= base || digit > int32(math.Floor(float64(maxInt-i)/float64(w))) {
				//return "", errors.New("overflow: input needs wider integers to process")
			}

			i += digit * int32(w)

			var t int32

			if int32(k) <= bias {
				t = tMin
			} else {
				if int32(k) >= bias+tMax {
					t = tMax
				} else {
					t = int32(k) - bias
				}
			}

			if digit < t {
				break
			}

			var baseMinusT = base - t
			if w > int(math.Floor(float64(maxInt)/float64(baseMinusT))) {
				//return "", errors.New("overflow: input needs wider integers to process")
			}

			w *= int(baseMinusT)

		}

		var out = len(output) + 1
		bias = p.Adapt(i-oldi, int32(out), oldi == 0)

		// `i` was supposed to wrap around from `out` to `0`,
		// incrementing `n` each time, so we'll fix that now:
		if int32(math.Floor(float64(i)/float64(out))) > maxInt-n {
			//return "", errors.New("overflow: input needs wider integers to process")
		}

		n += int32(math.Floor(float64(i) / float64(out)))
		i %= int32(out)

		// Insert `n` at position `i` of the output.
		output = Insert(output, int(i), n)
		i++
	}
	return string(output)
}

// 0 <= index <= len(a)
func Insert(a []int32, index int, value int32) []int32 {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}
