prod
====

[![Build Status](https://travis-ci.org/kusabashira/prod.svg?branch=master)](https://travis-ci.org/kusabashira/prod)

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

Options
-------

### --help

Display the usage and exit.

### --version

Output the version of prod and exit.

### -s, --separator=STRING

Use STRING instead of `\t` to separate columns.

```
$ prod <(seq 0 1) <(seq 0 1) <(seq 0 1)
0	0	0
0	0	1
0	1	0
0	1	1
1	0	0
1	0	1
1	1	0
1	1	1

$ prod --separator=, <(seq 0 1) <(seq 0 1) <(seq 0 1)
0,0,0
0,0,1
0,1,0
0,1,1
1,0,0
1,0,1
1,1,0
1,1,1
```

Other specification
-------------------

#### An empty file in `[FILE]...`

If an empty file is included in `[FILE]...`, prod outputs nothing.

```
$ prod <(printf "%s\n" A B) /dev/null <(seq 7)
```

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
