# CONUST

TODO: make string conversion the default and examine string before building, fractionals, doc, benchmark optimisation, ports

A utility to transform numbers into alphabetically sortable strings with the ability of reversing the transformation.

The there are two format options for the transformed strings, one uses the decimal (base(10)) representation of the number as its basis, the other the base(36) representation. The encoded values can hold numbers with up to 35 integral and 35 fractional base converted digits, which should be way more than necessary for everyday applications.

Trailing zeros of the integral part are encoded in one byte, and the leading zeros of the fractional part are encoded similarly. This means, that if your numbers typically contain lots of zeros, you probably should use the base(10) variant. If, however, your numbers are of arbitrary value, you might be able to save some space by using the base(36) variant as it is a more compact representation.

## Encoded Format Description

The zero value is transformed to simply the zero character ("5").

The first character of the output will be "6" for positive numbers, and "4" for negative numbers.

The length of the number is encoded in a base(36) digit followed by the number itself, omitting the trailing zeros.

In case of a negative numbers, all digits are swapped for their corresponding pair counting from the end of the digi alphabet. E.g. in base(10) a "2" becomes "7" and vice versa, in base(36) "1" becomes "y" and vice versa.
Furthermore, a "~" (tilde) is appended to the output for negative numbers, even integers.

The fractional part is separated from the integral part by "." (period) in case of positive numbers, and by a "~" (tilde) for negative numbers.
When encoding fractional numbers, you have to define the number of fractional digits to keep. The fractional input value will be rounded accordingly before processing.

If the number has a non zero fractional part, the output continues with the decimal separator character approppriate for the sign of the number (which is "." for positive and "~" for negative numbers).
The output will contain the number of leading zero digits of the fractional part in its new base encoded in a single digit.
Finally the remainig digits of the fractional part are output.