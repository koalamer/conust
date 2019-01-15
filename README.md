# CONUST

A utility to transform numbers into alphabetically sortable strings and vice versa.

The there are two format options for the transformed strings, one uses the decimal (Base(10)) representation of the number, the other it's Base(36) representation.
Trailing zeros of the integral part of each number are encoded in one byte, and the leading zeros of the fractional part are encoded similarly. This means, that if your numbers typically contain lots of zeros, you probably should use the Base(10) variant. If, however, your numbers are of arbitrary value, you might be able to save some space by using the Base(36) variant.

The encoded values can store numbers with up to 35 digits in their respective bases for both the integral and the fractional part. Both variants can hold the full range of an int32 or int64 values with one exception: since the minimum value of both int32 (-2 147 483 648) and int64 (-9 223 372 036 854 775 808) has an absolute value which is bigger by one than the maximum value that these types can represent, trying to encode them causes an error.
If there is a chance you may come across these minimum values, do check the input values before sending them to the encoder.  

## Encoded Format Description

The zero value is transformed to simply the zero character "5". This means that regardless wether the input was "0", "0.0","+0", "-0.0" or any other representation of zero, the output will simply be the zero character.

The first character of the output will be "9" for positive numbers, and "0" for negative numbers.

The length of the number is encoded in a Base(36) digit followed by the number itself, omitting the trailing zeros.

In case of a negative numbers, all digits are swapped for their corresponding pair counting from the end of the digi alphabet. E.g. in Base(10) a "2" becomes "7" and vice vera, in Base(36) "1" becomes "y" and vice versa. Furthermore, a "~" (tilde character) is appended to the output.

The fractional part is separated from the integral part by "." (period) in case of positive numbers, and by a "~" (tilde) for negative ones.
When encoding fractional numbers, you have to define the number of decimal digits to keep. The input value will be rounded accordingly before processing.
If the number has a non zero fractional part, the output continues with the decimal separator character approppriate for the sign of the number (which is "." for positive and "~" for negative numbers).
The output will contain the number of leading zero digits of the fractional part in its new base encoded in a single digit. After that, the digits of the fractional part are output the same way as was done with the integral part.