Gomig
=====

Synchronize data between databases, currently only supports syncing from
MySQL to PostgreSQL as that's what I need at the moment. The
architecture is loosely based on
[py-mysql2pgsql](https://github.com/philipsoutham/py-mysql2pgsql/) and
I've tried to keep it extensible, so with some help it should support
Postgres to MySQL, MySQL to Mysql, Oracle to Postgres et cetera.  Pull
requests welcome.

The default config file is called "config.yml", format is YAML,
and the parameters are a strict superset of
[py-mysql2pgsql](https://github.com/philipsoutham/py-mysql2pgsql/)'s
configuration format. If no config file is present on the first run, the
sample config file will be installed in its place.

Features
========
- Uses postgres' **COPY FROM** support for fast data transfers
- Define projections (views) in the source database so that they match a
  (reduced) form of tables in the destination database. **Gomig** will
  sync the data.
- Can execute SQL directly on the destination server or output to a
  file, just like
  [py-mysql2pgsql](https://github.com/philipsoutham/py-mysql2pgsql/).
- Will ROLLBACK when something goes wrong, leaving the destination
  database intact. The source database is never INSERT/UPDATE/DELETE'ed,
  only views or projection tables are created on request, they can be
  safely dropped should they somehow survive culling.

Running
=======
If you have the go runtime installed, then it's as easy as:

```bash
# get it!
$ go get github.com/aktau/gomig

# run it (assuming $GOPATH is in your $PATH, as I assume it would be for
# most Go devs, otherwise you'll have to cd to $GOPATH/bin)

# check out the options
$ gomig --help

# generate a config file, edit it, then run
$ gomig generate-config
$ edit config.yml
$ gomig
# alternatively you can explicitly supply a config file:
$ gomig -f config.yml
```

To update to the newest version later, you can just do:

```bash
$ go get -u github.com/aktau/gomig
```

Of course, with Go it's quite easy to just build binaries so that
everyone can run it without needing anything but the standard C runtime
(which is basically installed everywhere). But I'm not doing that just
yet. I might do that if there's interest though.

Build requirements
==================
- Go (>=) 1.2 (uses positional notation in fmt.Sprintf)

Todo
====
- Convert more datatypes, and do it more accurately. **Gomig** only
  handles varchar, text, blob (binary), boolean, integer and float at
  the moment. Dates/timestamps are in progress.
- Possibly faster data migration with goroutines, as explained in [this
  article](http://www.acloudtree.com/how-to-shove-data-into-postgres-using-goroutinesgophers-and-golang/).
  Would need to make things quite a bit more threadsafe for that though,
  or keep the goroutines internal to the bulk methods...

Screenshot
==========

Everybody loves screenshots, even of console apps.

![Fancy screenshot of gomig in action](http://aktau.github.io/gomig/images/screen-0.4.0-1.png)

License (BSD)
======================

Copyright (c) 2013, Nicolas Hillegeer
All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice, this
  list of conditions and the following disclaimer in the documentation and/or
  other materials provided with the distribution.

* Neither the name of the {organization} nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
