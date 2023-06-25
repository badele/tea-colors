For converting your terminal output use the unbuffer command (in expect package)


ex:

```console
echo "the output comment" > ls.ans
unbuffer ls -alh  >> samples/ls.ans
```