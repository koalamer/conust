# CONUST

A utility to transform numbers into alphabetically sortable strings with the ability of reversing the transformation. It is meant to be used when text tokens and numbers are stored both as strings and you need proper sorting on them using simple string sorting.

The input for the encoding must be a numeric string. It need not be integer, floating point numbers are accepted as well. The input can be in a base between 2 and 36. If the input has a base higher than 10, and contains letters, those must be lower cased.

The encoded version can be 1 - 3 characters longer than the original, but on the other hand the transformation only keeps the significant section of the number, removing all trailing and heading zeros, thereby possibly saving some space.

Beside the simple Encode and Decode functions that deal with individual numeric strings, there is the EncodeMixedText convenience function that scans the input for subsequent decimal characters and creates an output where these are encoded by Encode and surrounded by spaces.

## Conust for other languages

Currently there is only this Go version, but the converted format is simple to implement. See the next section if you would like to give it a try. (I might do ports myself later.)

## Encoded Format Description

Encoding an empty string results in an empty string.

For non empty input all trailing and heading zeros are ignored, and the first digit of the encoded number X will be:

- "7" if X >= 1
- "6" if 1 > X > 0
- "5" if X = 0, and threre are no more characters in this case
- "4" if 0 > X > -1
- "3" if -1 >= X

This is followed by the magnitude value of the significant part of the number, which can occupy more than one digit. The value of the magnitude is

- the number of integer digits when X >= 1 or X <= -1
- the number of leading zeros after the decimal point when X < 1 and X > -1 but X != 0

The value of the magnitude (M) is stored in a series of digits, each digit adding a maximum of 34 to the overall value of the magnitude:

- if 0 <= M <= 34 this value is stored in one digit
- if M > 34, then a digit vith the value of 35 is stored, and the encoding is repeated for the value M - 34

For numbers with the sign digit of

- "7" or "4" the magnitude digits are normal base 36 digits.
- "6" and "3" the digits are value inverted: instead of X there will be the digit 35 - X

After the magnitude come the significant digits of the original number, omitting the decimal point is there is any. The digits are treated as base 36 digits and are encoded:

- as normal digits if the number is positive, which basically means thet the digits are copied from the input
- as inverted digits if the number is negative, meaning that instead of digit X, the digit 35 - X is stored

Finally if the number is negative it is terminated by a "~" (tilde) character

## Conversion Examples

You can find conversion test data in the test files, but to showcase a few scenarios (in which by inverted I mean each digit X being substituted with digit 35 - X):

| input | encoded version | sing byte | magnitude | significant digits |
|---|---|---|---|---|
| 120000000000000000000000000000000 | 7z412 | 7 (x>=1) | z4 (34+4=38) | 12 |
| 1200 |7412 | 7 (x>=1) | 4 | 12 |
| 12 |7212 | 7 (x>=1) | 2 | 12 |
| 1.2 |7112 | 7 (x>=1) | 1 | 12 |
| 0.12 |6z12 | 6 (1>x>0) | z (0 inverted) | 12 |
| 0.0012 |6x12 | 6 (1>x>0) | x (2 inverted) | 12 |
| 0.0000000000000000000000000000000000012 | 60y12 | 6 (1>x>0) | 0y (z1 inverted) | 12 |
| 0 | 5 | 5 (x=0) |  |  |
| -0.0000000000000000000000000000000000012 | 4z1yx~ | 4 (0>x>-1) | z1 (34+1=35) | yx (12 inverted) |
| -0.0012 | 42yx~ | 4 (0>x>-1) | 2 | yx (12 inverted) |
| -0.12 | 40yx~ | 4 (0>x>-1) | 0 | yx (12 inverted) |
| -1.2 | 3yyx~ | 3 (-1>x) | y (1 inverted) | yx (12 inverted) |
| -12 | 3xyx~ | 3 (-1>x) | x (2 inverted) | yx (12 inverted) |
| -1200 | 3vyx~ | 3 (-1>x) | v (4 inverted) | yx (12 inverted) |
| -12000000000000000000000000000000000000 | 30vyx~ | 3 (-1 > x) | 0v (z4 inverted) | yx (12 inverted) |
