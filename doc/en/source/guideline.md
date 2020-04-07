# Contribution and Coding Guideline

## FogFlow Contribution Guide

This document describes the guidelines to contribute to FogFlow. If you are
planning to contribute to the code you should read this document and get familiar with its content.

## General principles

* FogFlow uses GO programing language (although other  tools like test tools and other script can be written in python , java, bash ).
* Efficient code (i.e. the one that achieves better performance) is preferred upon inefficient code. Simple code
  (i.e. cleaner and shorter) is preferred upon complex code. Big savings in efficiency with small penalty in
  simplicity are allowed. Big savings in simplicity with small penalty in efficiency are also allowed.
* Code contributed to FogFlow must follow the [code style guidelines](#code-style-guidelines)
  in order to set a common programming style for all developers working in the code.

Note that contribution workflows themselves (e.g. pull requests, etc.) are described in another document
([FIWARE Development Guidelines](https://forge.fiware.org/plugins/mediawiki/wiki/fiware/index.php/Developer_Guidelines)).

## Pull Request protocol

As explained in [FIWARE Development Guidelines](https://forge.fiware.org/plugins/mediawiki/wiki/fiware/index.php/Developer_Guidelines)
contributions are done using a pull request (PR). The detailed "protocol" used in such PR a is described below:

* Direct commits to master branch (even single-line modifications) are not allowed. Every modification has to come as a PR.
* In case the PR is implementing/fixing a numbered issue, the issue number has to be referenced in the body of the PR at creation time.
* Anybody is welcome to provide comments to the PR (either direct comments or using the review feature offered by Github).
* Use *code line comments* instead of *general comments*, for traceability reasons (see comments lifecycle below).
* Comments lifecycle
  * Comment is created, initiating a *comment thread*.
  * New comments can be added as responses to the original one, starting a discussion.
  * After discussion, the comment thread ends in one of the following ways:
    * `Fixed in <commit hash>` in case the discussion involves a fix in the PR branch (which commit hash is
       included as reference).
    * `NTC`, if finally nothing needs to be done (NTC = Nothing To Change).
    
 * PR can be merged when the following conditions are met:
    * All comment threads are closed.
    * All the participants in the discussion have provided a `LGTM` general comment (LGTM = Looks good to me)
 * Self-merging is not allowed (except in rare and justified circumstances).

Some additional remarks to take into account when contributing with new PRs:

* PR must include not only code contributions, but their corresponding pieces of documentation (new or modifications to existing one) and tests.
* PR modifications must pass full regression based on existing test (unit, functional, memory, e2e) in addition to whichever new test added due to the new functionality.
* PR should be of an appropriated size that makes review achievable. Too large PRs could be closed with a "please, redo the work in smaller pieces" without any further discussing.


## Code style guidelines

Note that, currently not all FogFlow's existing code base conforms to these rules. There are some parts of the code that were
written before the guidelines were established. However, all new code contributions MUST follow these rules and, eventually, old code will be modified to conform to the guidelines.

### ‘MUST follow’ rules

#### M1 (Headers Files Inclusion):

*Rule*: All header or source files MUST include all the header files it needs AND NO OTHER header files. They MUST
NOT depend on inclusions of other header files. Also, all header and source files MUST NOT include any header files it
does not need itself.

*Rationale*: each file should not depend on the inclusions other files have/don’t have. Also, if a header file
includes more files than it needs, its ‘clients’ has no other choice than to include those ‘extra’ files as
well. This sometimes leads to conflicts and must be avoided. In addition, it increases the compilation time.

*How to check*: manually

#### M2 (Copyright header)

*Rule*: Every file, source code or not, MUST have a copyright header:

For Golang files:

```
/*
*
* Copyright 20xx The FogFlow Authors.
*
* This file is part of FogFlow.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
*
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software,
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/
```
For Python, bash script etc.:

```
# Copyright 20XX FogFlow Authors.

# This file is part of FogFlow.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# For Python, bash script  etc.:

# Author: <the author>
```
*Rationale*: to have a homogenous copyright header for all files.

*How to check*: manually

#### M3 (Function header)

*Rule*: All functions MUST have a header, which SHOULD have a short description of what the function does, a
descriptive list of its parameters, and its return values.

Example:

```
/* ****************************************************************************
*
* parseUrl - parse a URL and return its pieces
*
*  [ Short description if necessary ]
*
* PARAMETERS
*   - url         The URL to be examined
*   - host        To output the HOST of the URL
*   - port        To output the PORT of the URL
*   - path        To output the PATH of the URL
*   - protocol    To output the PROTOCOL of the URL
*
* RETURN VALUE
*   parseUrl returns TRUE on successful operation, FALSE otherwise
*
* NOTE
*   About the components in a URL: according to
*   https://tools.ietf.org/html/rfc3986#section-3,
*   the scheme component is mandatory, i.e. the 'http://' or 'https://' must
*   be present, otherwise the URL is invalid.
*/
```

*Rationale*: the code is simply easier to read when prepared like this

*How to check*: manually

#### M4 (Indent)

*Rule*: Use only spaces (i.e. no tabs), and indent TWO spaces at a time.

*Rationale*: two whitespaces are enough. It does not makes the lines too long

*How to check*: manually

#### M5 (Variable declaration):

*Rule*: Each declared variable MUST go on a separate line:

```
var  i  int;
var  j  int;
```

The following usage MUST be avoided:

```
var  i, j, k int;
```

*Rationale*: easier to read.

*How to check*: manually

#### M6 (Naming conventions):

*Rule*: the following naming conventions apply:

* A name must begin with a letter, and can have any number of additional letters and numbers.
* A function name cannot start with a number.
* A function name cannot contain spaces.
* If the functions with names that start with an uppercase letter will be exported to other packages. If the function name starts with a lowercase letter, it won't be exported to other packages, but you can call this function within the same package.
* If a function name consists of multiple words, use camel case to represent such names, for example: empName, empAddress, etc.
* function names are case-sensitive (car, Car and CAR are three different variables).

*Rationale*: this rule makes it easy to understand.

*How to check*: manually

#### M7 (Use gofmt before commit for indentation and other formatting):

*Rule*: gofmt -r '(a) -> a' -w FileName

* Code before applying gofmt

```
package main
          import "fmt"
// this is demo to format code
            // with gofmt command
 var a int=2;
             var b int=5;
                            var c string= `hello world`;
       func print(){
                   fmt.Println("Value for a,b and c is : ");
                        fmt.Println(a);
                                 fmt.Println((b));
                                         fmt.Println(c);
                         }
```
* Code after applying rule

```
package main
 
import "fmt"
 
// this is demo to format code
// with gofmt command
var a int = 2
var b int = 5
var c string = `hello world`
 
func print() {
        fmt.Println("Value for a,b and c is : ")
        fmt.Println(a)
        fmt.Println((b))
        fmt.Println(c)
}
```
*Note use gofmt /path/to/package for package formating.
*Rationale*: This will reformat the code and updates the file.
*How to check*: manually

#### M8 (Command & operators separation):

*Rule*: operators (+, *, =, == etc) are followed and preceded by ONE space. Commas are followed by ONE space.

```
FogFunction(va`r1, var2, var3) {
	if (var1 == var2) {
  		var2 = var3;
	}
}
```

not

```
FogFunction(var1,var2,var3) {
	if (var1==var2) {
  		var1=var3;
	}
}
```

*Rationale*: easier on the eye.

*How to check*: manually

### ‘MUST follow’ rules 

#### S1 (Error management):

*Rule*: Error returned in the second argument should be managed.

* Bad implementation
```
FogContextElement, _ := preprocess(UpdateContextElement)
```
* Good implementation

```
preprocessed, err := preprocess(bytes)
if err != nil {
  return Message{}, err
}
```

#### S2 (Error printing message):

*Rule*: An error string shall neither be capitalized nor end with a punctuation according to Golang standards.

* Bad implementation
```
if len(in) == 0 {
  return "", fmt.Errorf("Input is empty")
}
```
* Good implementation

```
if len(in) == 0 {
	return nil, errors.New("input is empty")
}
```
#### S3 (Avoid nesting):

*Rule*: avoid nesting while writing the code.

* Bad implementation

```
func FogLine(msg *Message, in string, ch chan string) {
    if !startWith(in, stringComment) {
        token, value := parseLine(in)
        if token != "" {
            f, contains := factory[string(token)]
            if !contains {
                ch <- "ok"
            } else {
                data := f(token, value)
                enrichMessage(msg, data)
                ch <- "ok"
            }
        } else {
            ch <- "ok"
            return
        }
    } else {
        ch <- "ok"
        return
    }
}
```
* Good implemetation

```
func FogLine(in []byte, ch chan interface{}) {
    // Filter empty lines and comment lines
    if len(in) == 0 || startWith(in, bytesComment) {
        ch <- nil
        return
    }

    token, value := parseLine(in)
    if token == nil {
        ch <- nil
        log.Warnf("Token name is empty on line %v", string(in))
        return
    }

    sToken := string(token)
    if f, contains := factory[sToken]; contains {
        ch <- f(sToken, value)
        return
    }

    log.Warnf("Token %v is not managed by the parser", string(in))
    ch <- nil
}
```

#### S4 (Preconditions)

*Rule*: we strongly recommend for functions to evaluate the parameters and if necessary return error, before starting to process. 

* Bad implementation
```
a, err := f1()
if err == nil {
    b, err := f2()
    if err == nil {
        return b, nil
    } else {
        return nil, err
    }
} else {
    return nil, err
}
```
* Good implementation

```
a, err := f1()
if err != nil {
    return nil, err
}
b, err := f2()
if err != nil {
    return nil, err
}
return b, nil
```

#### S5 (If condition)
*Rule*: Go have some improved version in if condition 

* Bad implementation in Golang

```
f, contains := array[index]
if contains {
    // Do something
}
```
* Good implementation

```
if f, contains := array[index]; contains {
    // Do something
}
```
#### S5 (Switch)
*Rule*: always use default with switch condition.

* Bad implementation
```
switch simpleToken.token {
case tokenTitle:
    msg.Title = value
case tokenAdep:
    msg.Adep = value
case tokenAltnz:
    msg.Alternate = value 
// Other cases
}
```

* Good implementation 

```
switch simpleToken.token {
case tokenTitle:
    msg.Title = value
case tokenAdep:
    msg.Adep = value
case tokenAltnz:
    msg.Alternate = value
// Other cases    
default:
    log.Errorf("unexpected token type %v", simpleToken.token)
    return Message{}, fmt.Errorf("unexpected token type %v", simpleToken.token)
}
```
#### S5 (Constants management)

*Rule*:Constant value should be managed by ADEXP and ICAO message

* Bad implementation
```
const (
    AdexpType = 0 // TODO constant
    IcaoType  = 1
)
```
* Good implementation 

```
const (
    AdexpType = iota
    IcaoType 
)
```

