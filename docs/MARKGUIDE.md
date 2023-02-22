# Nemomark Guide

Nemomark is special syntax for help writing rich document.

# First step

Nemomark is made up of "functions". Basic form of function is like this:

``` 
$[func-name(argument1=value1) Content]
```

For example, if you want to bold a text, you can use function like this:

```
$[bold Hello] New World!
```

When text is contain Nemomark's function syntax string, Nemomark cannot display the text correctly.

To solve problem, you can use "Ignore" syntax.

``` 
This is string contains dollar sign`$` and square brackets`[]`

> This is string that contains dollar sign$ and square brackets[]
```

For example, this kind of text is not be displayed correctly.

```
func foo(bar string[]) string {
    //This is sample code to explain somethings.
}

> func foo(bar string
```

```
`func foo(bar string[]) string {
    //This is sample code to explain somethings.
}`

> func foo(bar string[]) string {
>   //This is sample code to explain somethings.
> }
```

# Functions

#### Bold

``` $[bold Text] ```

#### Italic

``` $[italic Text] ```

#### Strikehrough

``` $[cancel Text] ```

#### Underline

``` $[underline Text] ```

#### Link

``` $[link(url=example.com) Text] ```

#### Image

``` $[image(url=imagesrc) Alt Text] ```

To use your own image,

1. Put your image file to ``"/post/res"``.
2. Link your image like this:
   ``` $[image(url=./res/yourimage.png)] ```

#### Code

**WE RECOMMEND USING "IGNORE SYNTAX"**

``` $[code `func yourcode(foo strings[]){ //This is Sample Code }`] ```
