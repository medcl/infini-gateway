# nlpt_detect
Uses a go wrapper around the [cld2](https://code.google.com/p/cld2/) project, which is a c++ Naïve Bayesian classifier open-sourced by Google Chrome team. It "probabilistically detects over 80 languages in Unicode UTF-8 text, either plain text or HTML/XML."

## Get it

```sh
go get github.com/jbowles/cld2_nlpt
```

## Use it

```go
	english := "This is an english sentence"
	got := Detect(english, "name", len(english), 3, 3, 3)
  // -> ENGLISH
```

So far there are only 5 functions that employ only 3 different ways to use CLD2.

## Dependencies
The go wrapper: [cld2_nlpt](https://github.com/jbowles/cld2_nlpt), which is itself dependent on cld2 itself.

## nlpt
natural language processing toolkit
