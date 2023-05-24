# EBNF
[![GoDoc](https://godoc.org/github.com/mdhender/eebnf?status.svg)](http://godoc.org/github.com/mdhender/ebnf)

EBNF is an implementation of
[Niklaus Wirth](https://en.wikipedia.org/wiki/Niklaus_Wirth)'s
[A+DS=P](https://en.wikipedia.org/wiki/Algorithms_%2B_Data_Structures_%3D_Programs)'s parser in Go.

## Sources
Forked from
[golang.org/x/exp/ebnf](https://pkg.go.dev/golang.org/x/exp/ebnf).
That source is licensed under the following conditions:

    Copyright (c) 2009 The Go Authors. All rights reserved.
    
    Redistribution and use in source and binary forms, with or without
    modification, are permitted provided that the following conditions
    are met:
    
      * Redistributions of source code must retain the above copyright
        notice, this list of conditions and the following disclaimer.
      * Redistributions in binary form must reproduce the above
        copyright notice, this list of conditions and the following
        disclaimer in the documentation and/or other materials provided
        with the distribution.
      * Neither the name of Google Inc. nor the names of its contributors
        may be used to endorse or promote products derived from this
        software without specific prior written permission.
    
    THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
    "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
    LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
    A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
    OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
    SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
    LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
    DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
    THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
    (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
    OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

https://cs.opensource.google/go/x/exp/+/master:ebnf/

The internal builder is an update of
[Shivam Mamgain](https://github.com/shivamMg)'s
[recursive descent parser builder](https://github.com/shivamMg/rd).
His
[pretty printer](https://github.com/shivamMg/ppds)
is used to print out trees for debugging.