env-awsps
---
set from aws parameter store to environment variable.

Usage
---
```
$(env-awsps --region ap-northeast-1 --prefix app.web. --rm-prefix)
```

Command Line Option
---
- -h,--help: display help
- --version: display version and revision
- --prefix string: specify prefix of key
- --rm-prefix: remove prefix of key
- --region string: specify aws region

Example
---
```
(set value to parameter store)
app.web.user = XX1
app.web.password = YY1
app.db.user = XX2
app.db.password = YY2
```

```
(get value from parameter store)
$ env-awsps
export APP_WEB_USER=XX1
export APP_WEB_PASSWORD=YY1
export APP_DB_USER=XX2
export APP_DB_PASSWORD=YY2

$ env-awsps --prefix app.web.
export APP_WEB_USER=XX1
export APP_WEB_PASSWORD=YY1

$ env-awsps --prefix app.web. --rm-prefix
export USER=XX1
export PASSWORD=YY1
```

LICENSE
---
MIT
