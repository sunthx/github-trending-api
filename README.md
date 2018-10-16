# github-trending-api-server

提供了两个接口：`trending` 和 `contribution`

## trending

返回当前 Github Trending 的内容

```bash
curl 192.144.166.24:8090/trending
```

```json
{
    "Repositories": [
        {
            "name": "/leonardomso/33-js-concepts",
            "description": "<g-emoji class=\"g-emoji\" alias=\"scroll\" fallback-src=\"https://assets-cdn.github.com/images/icons/emoji/unicode/1f4dc.png\">📜</g-emoji> 33 concepts every JavaScript developer should know.",
            "url": "https://github.com/leonardomso/33-js-concepts",
            "star": "2,311",
            "fork": "93",
            "lang": "JavaScript"
        },
        {
            "name": "/open-source-for-science/TensorFlow-Course",
            "description": "Simple and ready-to-use tutorials for TensorFlow",
            "url": "https://github.com/open-source-for-science/TensorFlow-Course",
            "star": "2,078",
            "fork": "102",
            "lang": "Python"
        },
        {
            "name": "/photoprism/photoprism",
            "description": "Personal photo management powered by Go and Google TensorFlow",
            "url": "https://github.com/photoprism/photoprism",
            "star": "1,329",
            "fork": "33",
            "lang": "Go"
        },
        {
            "name": "/Igglybuff/awesome-piracy",
            "description": "A curated list of awesome warez and piracy links",
            "url": "https://github.com/Igglybuff/awesome-piracy",
            "star": "1,628",
            "fork": "114",
            "lang": ""
        },
        {...}
    ]
}

```

## contribution

返回指定码农过去一年的贡献数

```bash
curl 192.144.166.24:8090/contributions?user=sunthx
```

```json
[
    {
        "count": 0,
        "date": "2017-10-15",
        "color": "#ebedf0"
    },
    {
        "count": 0,
        "date": "2017-10-16",
        "color": "#ebedf0"
    },
    {
        "count": 0,
        "date": "2017-10-17",
        "color": "#ebedf0"
    },
    {
        "count": 0,
        "date": "2017-10-18",
        "color": "#ebedf0"
    },
    {...}
]
```

## TODO

### trending 

- [ ] 加入语言，时间等参数
- [ ] Rank 榜单

### contribution

- [ ] 统计
