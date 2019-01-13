# CONUST
A utility to transform numbers into alphabetically sortable strings and vice versa.

# Format Description

The there are two format options for the transformed strings, the decimal format, which uses format requires the definition of a set of characters that will be used as the digits, three special characters that will be used to signal the sign of the number and two decimal separators characters that separate the integral and fractional parts of the numbers (one for the positive and one for the negative domain).
There is a default for these parameters that should work well across different programming languages, character encodings, database systems etc., but you can override these defaults with any set of characters that follow the required rules.
For decimal numbers, a fractional precision needs to specified: the number of fractional digits that are to be kept. There is no default value for this.

The default set of digits is "0123456789abcdefghijklmnopqrstuvwxyz" which is a total of 36 characters and thus the transformed numbers will be a sort of Base(36) representation of the original number. By changing the number of possible digits, you change the base of the transformation and thereby the numeric range that can be described by it as well as the compactness of the transformed string (more on that later).
When you specify a custom digit set, you need to do it in the form of a series of characters in ascending order, as demonstrated by the default.

The decimal separator for positive numbers must be a character that is smaller than any of the digit characters. The default value is ".".
The decimal separator for negative numbers must be a character that is greater than any of the digit characters. The default value is "~".

The sign for negative characters must be smaller than the zero character and the zero character must be smaller than the positive sign character. The default values are: "N" as the negative sign, "O" (capital letter o) as the zero character and "P" as the positive sign.

The transformation goes as follows (examples assume the default parameters):

The zero value is transformed to simply the zero character ("O"). This means that regardless wether the input was "0", "0.0","+0", "-0.0" or any other representation of zero, the output will simply be the zero character.

The first character of the output will be the positive sign character ("P") for positive numbers, and the negative sign character ("N") for negative numbers.

The number itself will be normalized by removing the sign (as it has already been accounted for), the leading zeroes of the integer part and the trailing zeroes after the fractional part (if any). If the number is a fractional number with an integer part of 0, one 0 is kept as the integer part. if there is a zero fractional part, the decimal separator is removed as well. For example "-001230.045600" will become "1230.0456", "+000.0560" will become "0.056", "03.00" will become "3".

Positive and negative numbers will use the series of digits differently: positive numbers will use them as they were specified: in ascending order, but negative numbers will use them in reverse order.
For example for positive numbers the 0 digit will be "0", the digit 10 will be "A", the digit 35 will be "Z". For negative numbers the 0 digit will be "Z", the digit 10 will be "P" and the digit 35 will be "0".

The integer part is transformed into its Base(the_number_of_possible_digits) representation of itself (Base(36) with the default parameters).

The second character of the output output will contain a single digit that describes how long the number is in its new base. If this new length is represented by L, the output character will be the digit L-1. (Since L is least 1, the -1 offset is there to utilize all possible digits for this position in the output, including the 0 digit.)
For example for a number that is 5 digits long in it's new base, the second character of the output will be "4".

The trailing zeroes of the new base number are trimmed, and the remaining part is output using the specified digit characters.

If the number has a fractional part, the output continues with the decimal separator character approppriate for the sign of the number (which is "," for positive and "~" for negative numbers).
The fractional part is multiplied to move the number of significant digits into the integer range, the rest is rounded off.
The output will contain the number of leading zero digits of the fractional part in its new base encoded in a single digit with the reverse digit alphabet.
Then the fractional part is then output the same way as the integer part has been, but without the sign.

