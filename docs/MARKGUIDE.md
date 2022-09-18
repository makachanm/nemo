# Nemomark Guide
Nemomark is speical syntax for help writing rich document.

# First step 
Nemomark is basically made up with "functions". Basic form of function is like this:

``` 
$[func-name(argument1=value1) Content]
```

For example, if you want to bold a text, you can use function like this:
```
$[bold Hello] New World!
```

When text is contain Nemomark's function syntax string, Nemomark cannot do a corect display of text.

To solve problem like this, you can use "Ignore" syntax.
``` 
This is string that contains dollar sign`$` and square brackets`[]`

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

## Functions 

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

If you want use your own image, 
 1. Put your image file to "/post/res".
 2. Link your image like this: 
 ``` $[image(url=./res/yourimage.png)] ```

#### Code 
**RECOMMEND TO USE "IGNORE SYNTAX"**

``` $[code `func yourcode(foo strings[]){ //This is Sample Code }`] ```
