# Welcome To Glob!

## Introduction

_Glob_ is a Golang package that adds support for UNIX shell-like pattern matching for strings, commonly known as 'globbing'.

It is released under the 3-clause New BSD license. See [LICENSE.md](LICENSE.md) for details.

```golang
import glob "github.com/ganbarodigital/go_glob"

g := NewGlob("*.go")
fmt.Sprintf(g.Match("parser.go")) // prints true
```

## Table of Contents <!-- omit in toc -->

- [Introduction](#introduction)
- [Why Use Glob?](#why-use-glob)
  - [Who Is Glob For?](#who-is-glob-for)
  - [What Are The Differences Between This Package And The Golang Filepath.Match()?](#what-are-the-differences-between-this-package-and-the-golang-filepathmatch)
- [How Does It Work?](#how-does-it-work)
  - [Getting Started](#getting-started)
  - [What Is A Glob Pattern?](#what-is-a-glob-pattern)
  - [What Does A Glob Pattern Look Like?](#what-does-a-glob-pattern-look-like)
  - [What About Extended Globbing, Globstars, and GLOB_IGNORE?](#what-about-extended-globbing-globstars-and-glob_ignore)
  - [What Happens When A Match Method Is Called?](#what-happens-when-a-match-method-is-called)
  - [How Are Errors Handled?](#how-are-errors-handled)
- [What Do I Do If I Find A Valid Pattern That Glob Errors On / Returns The Wrong Result For?](#what-do-i-do-if-i-find-a-valid-pattern-that-glob-errors-on--returns-the-wrong-result-for)
- [Creating A Glob](#creating-a-glob)
  - [NewGlob()](#newglob)
- [Match Methods](#match-methods)
  - [Match()](#match)
  - [MatchShortestPrefix()](#matchshortestprefix)
  - [MatchLongestPrefix()](#matchlongestprefix)
  - [MatchShortestSuffix()](#matchshortestsuffix)
  - [MatchLongestSuffix()](#matchlongestsuffix)
- [Other Methods](#other-methods)
  - [Pattern()](#pattern)

## Why Use Glob?

Golang already has `filepath.Match()`. Why do we need another globbing package?

### Who Is Glob For?

_Glob_ is for anyone who wants globbing support to be as close to UNIX shell behaviour as possible.

We've built Glob for our [Scriptish](https://github.com/ganbarodigital/go_scriptish) and [ShellExpand](https://github.com/ganbarodigital/go_shellexpand) projects, when we realised that Golang's `filepath.Match()` couldn't do everything we needed.

### What Are The Differences Between This Package And The Golang Filepath.Match()?

There's two important differences:

* longest match vs shortest match behaviour of the `*` wildcard
* being able to match prefix or suffix vs matching whole string

In UNIX shell scripts, you can choose whether `*` matches as few characters as possible, or whether it matches as many characters as possible.

```bash
echo ${PARAM1#12*4}  # matches as few as possible
echo ${PARAM1##12*4} # matches as many as possible
```

Golang's `filepath.Match()` always matches as few characters as possible. It isn't possible to tailor its behaviour to suit.

Also, Golang's `filepath.Match()` only supports matching the whole string. While you can use the `*` wildcard to simulate prefix/suffix matching, `filepath.Match()` can't tell you how long the matching prefix or suffix is.

If you only need to know that your input string matches your globbing pattern, then `filepath.Match()` will be faster and a better choice.

## How Does It Work?

### Getting Started

Import _Glob_ into your Golang code:

```golang

import glob "github.com/ganbarodigital/go_glob"
```

Create a `Glob` struct, by calling [glob.NewGlob()](#newglob) with your globbing pattern:

```golang
myGlob := glob.NewGlob("abc*.go")
```

Once you have your `Glob` struct, use its methods to glob your strings:

```golang
// equivalent of calling `filepath.Match()`
success, err := myGlob.Match(myInput)

// find matching prefix
pos, success, err := myGlob.MatchShortestPrefix(myInput)
if success {
    prefix := myInput[:pos]
}
```

### What Is A Glob Pattern?

A _glob pattern_ (or just _pattern_ for short) can be used as a filter, and/or as a search term.

* a filter when used with [glob.Match()](#match)),
* both a filter and a search term when used with the other [match methods](#match-methods)

Historically, it was used in UNIX shells to find a list of matching filenames using a simple set of wildcards. (This is known as __pathname expansion__ today.) As UNIX shells became more powerful, they added the ability to manipulate the contents of strings. Instead of inventing new syntax, they took the existing pathname expansion support, and reused (most!) of it against arbitrary strings too.

It's such an integral part of using UNIX systems that many UNIX services and daemons have added their own support for globbing over the years.

### What Does A Glob Pattern Look Like?

A glob pattern is made from:

* `?` is a wildcard, that matches exactly one character
* `*` is a wildcard, that matches zero or more characters. Sometimes it can be greedy (match as many characters as possible), and sometimes it can be ungreedy (match as few characters as possible). It all depends on which match method you are calling.
* `[...]` matches any one of the characters inside the `[` and `]`.
* `[^...]` matches any one of the characters that are _not_ inside the `[` and `]`
* `[lo-hi]` matches any one of the characters defined by the range `lo-hi`
* `\` escapes the following character. Use this to tell Glob to treat characters like `*` as a normal char and not as a wildcard.

Any other characters in the pattern are treated as a requirement to match exactly that character.

### What About Extended Globbing, Globstars, and GLOB_IGNORE?

_Extended globbing_ adds support for _pattern lists_ and _alternates_. It's not supported in this release. We'd like to support it in the future, but no promises!

_Globstars_ are the `**` and `**/` wildcards. They're used in _pathname expansion_ to match all files, all directories, and sub-directories. Because `Glob` currently only deals with arbitrary strings, it doesn't make sense to implement _globstar_ support atm.

`GLOB_IGNORE` is an environment variable used in _pathname expansion_ as a second filter against filepaths that have matched the globbing pattern. Because `Glob` currently only deals with arbitrary strings, it doesn't make sense to implement _GLOB_IGNORE_ support atm.

### What Happens When A Match Method Is Called?

This information is mostly to help you if you run into bugs in the _Glob_ package. Try not to rely on it to make your code work. A future version of _Glob_ may implement globbing in a different way.

Whenever you call any of the [match methods](#match-methods), here's what happens:

* we convert the glob pattern into a compiled Golang regex. The regex will be different for each of the match methods.
* we use the regex to discover if the pattern matches your input string
* where necessary, we do some additional work to find out which string slice index to return back to you

If we have already compiled a Golang regex for your glob and matcher method, we reuse it instead of compiling it again. This helps performance (for example) if you're globbing against a list of filenames - any situation where you'd be calling the same match method multiple times.

Golang's regex engine uses what's called leftmost-match semantics. Most of the time, that's exactly the behaviour you want ... unless you're after the shortest suffix that matches your pattern. That's where we have to do some additional processing of the regex result to find the shortest match of your pattern.

### How Are Errors Handled?

Errors can only occur:

* if you use a _pattern_ that is somehow invalid, for example `abc[`
* if you use a _pattern_ that isn't correctly understood (yet) by _Glob_

All of the [match methods](#match-methods) return an `error` back to you.

## What Do I Do If I Find A Valid Pattern That Glob Errors On / Returns The Wrong Result For?

We've got a comprehensive test suite, which is kept up to date. Even so, there could be glob patterns that should work, but don't.

* It could be that our glob2regex process doesn't correctly translate the pattern (e.g. missing escaping)
* It could be that the resulting regex doesn't behave the way the glob pattern does in a real UNIX shell

When you run into a problem, here's what to do:

1. create a small example `bash` shell script that demonstrates the correct behaviour
2. please [open an issue here on GitHub](https://github.com/ganbarodigital/go_glob)
3. add your example shell script, and details of what _Glob_ is doing, to the issue

We're aiming for 100% compatibility with UNIX shell globbing behaviour _when applied to arbitrary strings_.

We can't accept requests to make _Glob_ behave differently to how globbing works within a UNIX shell.

## Creating A Glob

### NewGlob()

To create a glob, call `glob.NewGlob()` with your _globbing pattern_:

```golang

myGlob := NewGlob(myPattern)
```

This gives you a `Glob` that you can reuse as many times as you want.

## Match Methods

Use one of the following match methods to perform the actual globbing.

### Match()

```golang
func (g *Glob) Match (input string) (bool, error)
```

`Match()` determines if the whole input string matches the given glob pattern.

Returns:

* `true` if the whole input string matches the Glob pattern, `false` otherwise
* an error if the given Glob pattern cannot be compiled into a regex

Example:

```golang
myGlob := NewGlob("*.go")
success, err := myGlob.Match("match.go")
```

### MatchShortestPrefix()

```golang
func (g *Glob) MatchShortestPrefix (input string) (int, bool, error)
```

`MatchShortestPrefix()` returns the prefix of input that matches the glob pattern. It treats '*' as matching the minimum number of characters.

Returns:

* the end of the prefix that matches the glob pattern, suitable for you to use in a string slice
* `true` if the pattern matched; `false` otherwise
* an error if the given Glob pattern cannot be compiled into a regex

Example:

```golang
input := "path/to/folder"

myGlob := NewGlob("*/")
pos, success, err := myGlob.MatchShortestPrefix(input)
if err != nil {
    // ... handle err first
}
if !success {
    // input string did not match pattern
}

// if we get here, the pattern did match
//
// in this example, the prefix will be 'path/'
prefix := input[:pos]
```

### MatchLongestPrefix()

```golang
func (g *Glob) MatchLongestPrefix (input string) (int, bool, error)
```

`MatchLongestPrefix()` returns the prefix of input that matches the glob pattern. It treats '*' as matching the maximum number of characters.

Returns:

* the end of the prefix that matches the glob pattern, suitable for you to use in a string slice
* `true` if the pattern matched; `false` otherwise
* an error if the given Glob pattern cannot be compiled into a regex

Example:

```golang
input := "path/to/folder"

myGlob := NewGlob("*/")
pos, success, err := myGlob.MatchLongestPrefix(input)
if err != nil {
    // ... handle err first
}
if !success {
    // input string did not match pattern
}

// if we get here, the pattern did match
//
// in this example, the prefix will be 'path/to/'
prefix := input[:pos]
```

### MatchShortestSuffix()

```golang
func (g *Glob) MatchShortestSuffix (input string) (int, bool, error)
```

`MatchShortestSuffix()` returns the suffix of input that matches the glob pattern. It treats '*' as matching the minimum number of characters.

Returns:

* the start of the suffix that matches the glob pattern, suitable for you to use in a string slice
* `true` if the pattern matched; `false` otherwise
* an error if the given Glob pattern cannot be compiled into a regex

__BE AWARE__ that the returned position _can be equal to `len(input)`. This happens when the pattern legitimately matches an empty suffix.

Example:

```golang
input := "path/to/folder"

myGlob := NewGlob("/*")
pos, success, err := myGlob.MatchShortestSuffix(input)
if err != nil {
    // ... handle err first
}
if !success {
    // input string did not match pattern
}

// if we get here, the pattern did match
//
// in this example, the suffix will be '/folder'
suffix := ""
if pos <len(input) {
    suffix := input[pos:]
}
```

### MatchLongestSuffix()

```golang
func (g *Glob) MatchLongestSuffix (input string) (int, bool, error)
```

`MatchLongestSuffix()` returns the suffix of input that matches the glob pattern. It treats '*' as matching the maximum number of characters.

Returns:

* the start of the suffix that matches the glob pattern, suitable for you to use in a string slice
* `true` if the pattern matched; `false` otherwise
* an error if the given Glob pattern cannot be compiled into a regex

__BE AWARE__ that the returned position _can be equal to `len(input)`. This happens when the pattern legitimately matches an empty suffix.

Example:

```golang
input := "path/to/folder"

myGlob := NewGlob("/*")
pos, success, err := myGlob.MatchLongestSuffix(input)
if err != nil {
    // ... handle err first
}
if !success {
    // input string did not match pattern
}

// if we get here, the pattern did match
//
// in this example, the suffix will be '/to/folder'
suffix := ""
if pos <len(input) {
    suffix := input[pos:]
}
```

## Other Methods

### Pattern()

Use `Pattern()` to get the original pattern out of a prepared `Glob` (e.g. for logging / debugging purposes):

```golang
myGlob := NewGlob("/*")
fmt.Printf("glob pattern is: %s\n", myGlob.Pattern())
```