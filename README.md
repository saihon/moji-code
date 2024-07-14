# moji-code

moji-code is a command line tool that can be displays character code.

# Installation

```
go install github.com/saihon/moji-code@latest
```

# Usage

* Outputs ASCII codes.
    ```
    $ moji-code
    0 U+0000 NULL
    1 U+0001 SOH
    2 U+0002 STX
    3 U+0003 ETX
    4 U+0004 EOT
    5 U+0005 ENQ
    6 U+0006 ACK
    .
    .
    .
    ```

* Output with details such as categorizes.
    ```
    $ moji-code -V
    0 U+0000 NULL        Control Common Null
    1 U+0001 SOH         Control Common Start Of Heading
    2 U+0002 STX         Control Common End Of Text
    3 U+0003 ETX         Control Common End Of Transmission
    4 U+0004 EOT         Control Common End Of Transmission
    5 U+0005 ENQ         Control Common Enquiry
    6 U+0006 ACK         Control Common Acknowledgement
    .
    .
    .
    ```

* Outputs a table corresponding to the character.
    ```
    $ moji-code abc
    97 U+0061 a
    98 U+0062 b
    99 U+0063 c
    ```

* Outputs a table corresponding to the decimal number.
    ```
    $ moji-code -d 97 98 99
    97 U+0061 a
    98 U+0062 b
    99 U+0063 c
    ```

* Outputs a table corresponding to the hexadecimal number.
    ```
    $ moji-code -x 0061 62 63
    97 U+0061 a
    98 U+0062 b
    99 U+0063 c
    ```

* Outputs a corresponding to the specified range.
    ```
    $ moji-code -r a c
    97 U+0061 a
    98 U+0062 b
    99 U+0063 c
    
    $ moji-code -r -d 97 99
    97 U+0061 a
    98 U+0062 b
    99 U+0063 c
    
    $ moji-code -r -x 61 63
    97 U+0061 a
    98 U+0062 b
    99 U+0063 c
    ```

* Categorize using the standard Golang package unicode.
    ```
    $ moji-code -V a ! 1 🙄
        97 U+0061 a           Graphic Letter Latin Lowercase ASCII-Hex-Digit Hex-Digit
        33 U+0021 !           Graphic Punctuation Common Pattern-Syntax STerm Sentence-Terminal Terminal-Punctuation
        49 U+0031 1           Graphic Digit Common ASCII-Hex-Digit Hex-Digit
    128580 U+1F644 🙄         Graphic Symbolic Common
    ```
