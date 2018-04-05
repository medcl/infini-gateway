
PROXY, a simple elasticsearch proxy written in Go.

# Features
- Enable auto handling elasticsearch failure, via disk based indexing queue(WIP)
- Enable Opt-in TLS/HTTPS protect
- Enable replicate indexing request to multi remote elasticsearch clusters(WIP)
- Enable load balancing(indexing and search request), algorithm configurable(WIP)

# How to use

- First, setup upstream config in the `proxy.yml`.

```
plugins:
- name: proxy
  enabled: true
  algorithm: duplicate
  disk_queue: true
  upstream:
  - name: primary
    enabled: true
    timeout: 60s
    elasticsearch:
      endpoint: http://localhost:9200
      index_prefix: gopa-
      username: elastic
      password: changeme
  - name: backup
    enabled: false
    timeout: 60s
    elasticsearch:
     endpoint: http://localhost:9201
     index_prefix: gopa-
     username: elastic
     password: changeme

```
- Start the PROXY.

```
➜  elasticsearch-proxy ✗ ./bin/proxy
___  ____ ____ _  _ _   _
|__] |__/ |  |  \/   \_/
|    |  \ |__| _/\_   |
[PROXY] An elasticsearch proxy written in golang.
0.1.0_SNAPSHOT,

[04-05 19:30:13] [INF] [instance.go:23] workspace: data/APP/nodes/0
[04-05 19:30:13] [INF] [api.go:147] api server listen at: http://0.0.0.0:2900

```

- Done! Now you are ready to rock with it.

```
➜  elasticsearch-proxy ✗ curl -XGET http://localhost:2900/
{
  "name" : "XZDZ8qc",
  "cluster_name" : "elasticsearch",
  "cluster_uuid" : "AIqV7VYGT9G13WgucUVu9g",
  "version" : {
    "number" : "6.2.3",
    "build_hash" : "c59ff00",
    "build_date" : "2018-03-13T10:06:29.741383Z",
    "build_snapshot" : false,
    "lucene_version" : "7.2.1",
    "minimum_wire_compatibility_version" : "5.6.0",
    "minimum_index_compatibility_version" : "5.0.0"
  },
  "tagline" : "You Know, for Search"
}
```

Have fun!


License
=======
Released under the [Apache License, Version 2.0](https://github.com/medcl/elasticsearch-proxy/blob/master/LICENSE) .

