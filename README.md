```azure
数据迁移工具
make build
```
```azure
配置样例
dsn0: &dsn0
  "user:pasword@(10.0.0.203:3306)/hezi?charset=utf8&parseTime=True&loc=Local"
dsn1: &dsn1
  "user:pasword@(10.0.0.203:3306)/log?charset=utf8&parseTime=True&loc=Local"
startConditionVal: &startConditionVal
  - 1
startCondition: &startCondition
  {"id > ?":*startConditionVal}
task:
    - {"dsn0":*dsn0,"t0":"user","dsnType0":"","dsn1":*dsn1,"t1":"","dsnType1":"",
"startFuncType":1,"startCondition": *startCondition,"nextFuncType":2,"nextKey":"id", "deleteKey":"id"
"endFuncType":2, "endKey":"id", "endVal":2
    }
```

```azure
执行
./dbMigrate -y=config.yam.
```