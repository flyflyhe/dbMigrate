```text
数据迁移工具
make build
```
```text
dsnType:mysql,sqlserver,postgresql
startFuncTypeCustom  = 1 //自定义启动规则 需配置startCondition
startFuncTypeDefault = 2 //取表的第一行数据

nextFuncTypeCustom = 1 //自定义迭代器规则 暂不支持 gorm不支持游标 
nextFuncTypeId     = 2 //通过主键id 迭代 主键非id 可配置 eg:nextKey:"uid"

endFuncTypeCustom   = 1 //结束规则自定义
endFuncTypeId       = 2 //id比较
endFuncTypeDatetime = 3 //datetime比较
```
```yaml
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

```text
执行
./dbMigrate -y=config.yaml
```