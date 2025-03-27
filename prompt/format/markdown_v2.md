# Syntax

*bold \*text*
_italic \*text_
__underline__
~strikethrough~
||spoiler||
[inline URL](http://www.example.com/)
`inline fixed-width code` 

```
pre-formatted fixed-width code block
```

```python
# pre-formatted fixed-width code block written in the Python programming language
```

>Block quotation started
>Block quotation continued
>Block quotation continued
>Block quotation continued
>The last line of the block quotation
**>The expandable block quotation started right after the previous block quotation
>It is separated from the previous block quotation by an empty bold entity
>Expandable block quotation continued
>Hidden by default part of the expandable block quotation started
>Expandable block quotation continued
>The last line of the expandable block quotation with the expandability mark||

## ВАЖНО!

- Any character with code between 1 and 126 inclusively can be escaped anywhere with a preceding '\' character, in which case it is treated as an ordinary character and not a part of the markup. This implies that '\' character usually must be escaped with a preceding '\' character.

- Inside pre and code entities, all '`' and '\' characters must be escaped with a preceding '\' character.

- Inside the (...) part of the inline link, all ')' and '\' must be escaped with a preceding '\' character.

- In all other places characters '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!' must be escaped with the preceding character '\'.

- In case of ambiguity between italic and underline entities __ is always greadily treated from left to right as beginning or end of an underline entity, so instead of ___italic underline___ use ___italic underline_**__, adding an empty bold entity as a separator.
