# gt 

A tiny but powerful Go internationalisation (i18n) library.

## Installation

```sh
$ go get github.com/melvinmt/gt
```

## Usage
```go
package main

import (
    "fmt"
    "github.com/melvinmt/gt"
)

var g = &gt.Build{
    Index: gt.Strings{
        "homepage-title": {
            "en":    "Hello World!",
            "es":    "¡Hola mundo!",
            "zh-CN": "你好世界!",
        },
        "homepage-welcome": {
            "en":    "Welcome to %s, %s.",
            "es-LA": "Bienvenidos a %s, %s.",
            "tr-TR": "%s, %s'ya hoşgeldiniz.",
            "nl":    "Welkom bij %s, %s",
        },
        "eating": {
            "en":    "%s#name came to %s#chain and ate a %s#hamburger.",
            "tr-TR": "%s#chain'ya geldi ve %s#hamburger yedi %s#name.",
        },
        "invoice": {
            "en":    "You need to pay %5.2f#amount dollars in %d#days days.",
            "pt-BR": "Você precisa pagar %5.2f#amount dólares em %d#days dias.",
        },
    },
    Origin: "en", // the language you'll be translating from
}

func main() {
    // Key based:
    g.SetTarget("es")
    s1 := g.T("homepage-title")
    fmt.Println(s1) // outputs: ¡Hola mundo!

    // String based:
    g.SetTarget("zh-CN")
    s2 := g.T("Hello World!")
    fmt.Println(s2) // outputs: 你好世界!

    // Parse arguments:
    g.SetTarget("nl")
    s3 := g.T("Welcome to %s, %s.", "Github", "Melvin")
    fmt.Println(s3) // outputs: Welkom bij Github, Melvin

    // As you can see in the previous example, you can use regular printf verbs
    // to parse arguments. The problem that you're facing here is that the order
    // of words is different in some languages:

    g.SetOrigin("es-LA")
    g.SetTarget("tr-TR")
    s4, err := g.Translate("Bienvenidos a %s, %s.", "Github", "Melvin")
    if err != nil {
        panic(err)
    }
    fmt.Println(s4)
    // This outputs: Github, Melvin'ya hoşgeldiniz. This is roughly translated as:
    // Welcome to Melvin, Github.  Which is NOT what you want. You can solve this with
    // tag notation.

    // Tag notation:
    g.SetOrigin("en")
    s5 := "%s#name came to %s#chain and ate a %s#hamburger."
    fmt.Println(g.T(s5, "Tom", "Burger King", "Whopper"))
    // Outputs: Burger King'ya geldi ve Whopper yedi Tom.
    // Which is in a different order, but correctly translated.

    // You can use any legal printf verb in combination with tags:
    g.SetOrigin("pt-BR")
    g.SetTarget("en")
    fmt.Printf(g.T("Você precisa pagar %5.2f#amount dólares em %d#days dias."), 499.95, 5)
    // Outputs: You need to pay 499.95 dollars in 5 days.
}
```

## Error handling

T() always returns strings. You can use the more verbose Translate() method to handle errors.

```go
key := "dashboard-notice"
s, err := g.Translate(key, "Melvin", "Alex")
if err != nil {
  // do something
}
```

## Edge cases

It's not recommended to pass duplicate anonymous printf verbs to **gt**, e.g. `"%s, %s, %d"`. It will work when the target strings will keep the arguments in order, but when one language requires to format the string as `%s %d %s`, **gt** will fail because it doesn't know which `%s` to swap. You can easily solve this by tagging duplicate verbs: `%s#tag1 %s#tag2 %d`.

Even when **gt** fails, it will try to return the origin string with formatted arguments. In this way, even when a translation has failed (or simply doesn't exist yet), you can at least present something to the end user.

```go
g := &gt.Build{
    Index: gt.Strings{
        "incomplete": {
            "en": "Sorry %s, this string is not translated yet!",
        },
    },
    Origin: "en",
}
g.SetTarget("es")
s, err := g.Translate("incomplete", "John")
if err != nil {
    fmt.Println(s) // outputs: Sorry John, this string is not translated yet!
}
```
When a key is used and the key doesn't exist and/or an origin language is not set, it can't find any origin strings and will simply return the key.
## Docs

http://godoc.org/github.com/melvinmt/gt

### *01/21/2013*: v1.0.1

- fixed some bugs in reverse lookup and cleaning tags

### *01/20/2013*: v1.0.0

- initial version of gt
- gt passes all tests 
- wrote the documentation
