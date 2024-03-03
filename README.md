# Motivation
Well it seems that xargs wont work well with parallelism. When i run this the output is awfull
```shell
git ls-files -z | xargs -P16 -0n1 git blame -w --line-porcelain | perl -n -e '/^author (.+)$/ && print "$1\n"' | sort -f | uniq -c | sort -nr
```
so i made this.. maybe i make it more complete to replace xargs in my workflow, but i dont know.


# TODO
- [] Should work well with the -n argument of xargs. The default behaviour should be 1 call and set all arguments to the call
- [] Take into account max len of the arguments and act accordingly