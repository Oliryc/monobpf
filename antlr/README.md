# https://medium.com/dailyjs/compiler-in-javascript-using-antlr-9ec53fd2780f

```
```

Had to download antlr-4.7.2, the jar for 4.7.1 was corrupted (use https://www.antlr.org/download.html)

```
curl https://raw.githubusercontent.com/antlr/grammars-v4/master/ecmascript/ECMAScript.g4 --output grammars/ECMAScript.g4
antlr4 -Dlanguage=JavaScript -lib grammars -o lib -visitor -Xexact-output-dir grammars/ECMAScript.g4
```

