# CONUST

A utility to transform numbers into alphabetically sortable strings with the ability of reversing the transformation. It is meant to be used when text tokens and numbers are stored both as strings and you need proper sorting on them.

The input for the encoding must be a numeric string. It need not be integer, floating point numbers are accepted as well. The input can be in a base bewtween 2 and 36. If the input has a base higher than 10, and contains letters, those must be lower cased.

The encoded version might be a few characters longer than the original, but on the other hand the transformation only keeps the significant section of the number, removing all trailing and heading zeros.

## Conust for other languages

Currently there is only this Go version, but the converted format is simple to implement. See the next section if you would like to give it a try. (I might do ports myself later.)

## Encoded Format Description

Encoding an empty string results in an epmty string.

For non empty input all trailing and heading zeros are ignored, and the first digit of the encoded number X will be:

- "7" if X >= 1 
- "6" if 1 > X > 0
- "5" if X = 0, and threre are no more characters in this case
- "4" if 0 > X > -1
- "3" if -1 >= X 

This is followed by the exponential of the significant part of the number, which can occupy more than one digit. The value of the exponent is

- if X >= 1  or X <= -1, the number of integer digits
- if X < 1 and X > -1 but X != 0, then the number of leading zeros after the decimal point is stored

The value of the exponent (E) is stored in a series of digits, each adding a maximum of 34 to the overall value of the exponent:

- if 0 <= E < 34 this value is stored in one digit
- if E > 34, then a digit vith the value of 35 is stored, and the encoding is repeated for E = E - 34

For numbers with the sign digits

- "7" or "4" the exponent digits are normal base 36 digits.
- "6" and "3" the digits are value reversed: instead of X there will be the digit 35 - X

After the exponential come the significant digits of the original number, omitting the decimal point is there is any. The digits are treated as base 36 digits and are encoded

- normally if the number is positive, which basically means thet the digits are copied from the input
- value reversed, meaning that  instead of a digit X, the digit 35 - X is stored

Finally if the number is negative it is terminated by a "~" (tilde) character

## Conversion Examples

You can find conversion test data in the test files, but to showcase a few scenarios:

| input | encoded version |
|---|---|
| 12000000000000000000000000000000000000 | 7z412 |
| 1200 |7412 |
| 12 |7212 |
| 1.2 |7112 |
| 0.12 |6z12 |
| 0.0012 |6x12 |
| 0.0000000000000000000000000000000000012 | 60y12 |
| 0 | 5 |
| -0.0000000000000000000000000000000000012 | 4z1yx~ |
| -0.0012 | 42yx~ |
| -0.12 | 40yx~ |
| -1.2 | 3yyx~ |
| -12 | 3xyx~ |
| -1200 | 3vyx~ |
| -12000000000000000000000000000000000000 | 30vyx~ |
