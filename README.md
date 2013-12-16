Gomig
=====

Synchronize data between databases, currently only supports syncing from
MySQL to PostgreSQL as that's what I need at the moment. The
architecture is loosely based on
[py-mysql2pgsql](https://github.com/philipsoutham/py-mysql2pgsql/) and
I've tried to keep it extensible, so with some help it should support
Postgres to MySQL, MySQL to Mysql, Oracle to Postgres et cetera.  Pull
requests welcome.

At the moment this small commandline app is not very friendly yet. The
config file is always called "config.yml", format is YAML, and the
parameters are a strict superset of
[py-mysql2pgsql](https://github.com/philipsoutham/py-mysql2pgsql/)'s
configuration format.

Features
========
- Uses postgres' **COPY FROM** support for fast data transfers
- Define projections (views) in the source database so that they match a
  (reduced) form of tables in the destination database. **gomig** will
  sync the data.
- Can execute SQL directly on the destination server or output to a
  file, just like
  [py-mysql2pgsql](https://github.com/philipsoutham/py-mysql2pgsql/).

Requirements
===========
- Go 1.2 (uses positional notation in fmt.Sprintf)

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
