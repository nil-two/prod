prod
====

Output direct product of lines of each files.

```
$ prod <(printf "%s\n" A B C) <(seq 2)
A	1
A	2
B	1
B	2
C	1
C	2
```

Usage
-----

```
$ prod [OPTION]... [FILE]...
Output direct product of lines of each files.

Options:
  -s, --separator=STRING    use STRING to separate columns (default: \t)
      --help                display this help text and exit
      --version             display version information and exit
```

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
