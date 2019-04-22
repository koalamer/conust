# CONUST

A utility to transform numbers into alphabetically sortable strings with the ability of reversing the transformation. It is meant to be used when text tokens and numbers are stored both as strings and you need proper sorting on them.

The input for the encoding must be a numeric string. It need not be integer, floating point numbers are accepted as well. The input can be in a base bewtween 2 and 36. If the input has a base higher than 10, and contains letters, those must be lower cased.

The encoded version might be 2 to 4 characters longer than the original, but on the other hand the transformation tries to save some space by ommitting the trailing zeros of the integral part of the number, and the leading zeros of the fractional part.
In case your numbers are such that they cannot benefit from this optimization, you might want to convert your numbers to a higher base before encoding, to make them shorter. There is a FloatConverter included in this library to help you with that, if you want to convert floating point numbers. For integers there is the strconv.FormatInt function to help you.

The way the numbers are transformed imposes limitations on what can be converted: the integral part of the number cannot be more than 35 digits long, and the fractional part cannot have more than 35 leading zeros. This still allows insanely large numbers, so it shouldn't be a problem.

## FloatConverter

The bundled FloatConverter is able to base convert float64 numbers. Integers are no problem, but fractional numbers might be represented in another base with a long (or even endless) series of fractional digits. For this reason you can specify what precision the converted value should retain. The default is to keep enough fractional digits in the new base to ensure the precision equivalent of 3 decimal digits (0.001 precision)

## Conust for other languages

Currently there is only this Go version, but the converted format is simple to implement. See the next section if you would like to give it a try. (I might do ports myself later.)

## Encoded Format Description

Encoding (or decoding) an empty string results in an epmty string.

All leading and trailing zeros are ignored on the input.

The zero value is transformed to simply the zero character ("5").

If the value is not zero, the first character of the output will be "6" for positive numbers, and "4" for negative numbers. This way the first character is always a number, which in the case of alphabetic ordering means that they will be sorted where numbers are normally sorted too. The rest of the encoding depends on the sign:

**For positive numbers**, after the sign digit ("6"):

The length of the integer part of the number is encoded in a single base(36) digit. This is followed by the integer itself, but without its trailing zeros. If the integer part is 0, then that one 0 IS present and the length is encoded as being 1.

The fractional part is separated from the integral part by "." (period).

The number of leading zeros of the fractional part (say X) is encoded in a single base(36) digit, but instead of X itself, 35 - X is used.

Then the fractioal digits following the trailing zeros are output.

**For negative numbers**, after the sign digit ("4"):

The technique is the same, but

- all digits are "reversed", meaning that instead of digit X you'll get digit 35 - X,
- instead of the decimal point you'll have a "~" (tilde) character,
- at the end of the output an extra "~" is added, if it is an integer it ends with two "~" characters.

## Conversion Examples

You can find conversion test data in the test files, but to showcase a few scenarios:

| input | encoded version |
|---|---|
| 0, +0, -0, 000, 0.0 | 5 |
| 12 | 6212 |
| 1200 | 6412 |
| -200 | 4wx~~ |
| 0.01 | 610.y1 |
| -0.01 | 4yz~yy~ |
| fcd200 | 66fcd2 |
| -5h32m.002d | 4uuiwxd~xxm~ |
