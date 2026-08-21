# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

灾备管理员恢复一份校验通过的牛群快照时，存储服务明确报了不可用，恢复数据也没有落库；系统却只留下 outcome=failed 的审计记录，接口最终仍返回成功，控制台随即把任务标成 completed。请先不要修改代码，定位这个相互矛盾的结果是怎样产生的，说明计划计算、数据写入、失败审计和最终返回值之间的完整传播链，并界定哪些输入路径会误报成功。

## 含 Bug 版本

- 仓库：VanceMichael/go-label-26
- 仓库地址：https://github.com/VanceMichael/go-label-26.git
- parent SHA：9111939b051643fdf8dcccdda0e8dbe98b1850f1

## 复现步骤

```bash
git clone -- https://github.com/VanceMichael/go-label-26.git bug-repro
cd bug-repro
git checkout --detach 9111939b051643fdf8dcccdda0e8dbe98b1850f1
go test ./internal/recovery -run ^TestRestoreReportsPersistenceFailureAfterAuditSucceeds$ -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/recovery -run ^TestRestoreReportsPersistenceFailureAfterAuditSucceeds$ -count=1
--- FAIL: TestRestoreReportsPersistenceFailureAfterAuditSucceeds (0.00s)
    coordinator_test.go:36: restore error = <nil>, want persistence failure
FAIL
FAIL	go-base/internal/recovery	0.035s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/recovery -run ^TestRestoreReportsPersistenceFailureAfterAuditSucceeds$ -count=1
--- FAIL: TestRestoreReportsPersistenceFailureAfterAuditSucceeds (0.00s)
    coordinator_test.go:36: restore error = <nil>, want persistence failure
FAIL
FAIL	go-base/internal/recovery	0.001s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

根因结论必须定位恢复协调器在持久化失败分支选择最终返回值的具体错误传播行为，串联计划和内存应用成功、Replace 失败、failed 审计成功以及上层误判完成的先后关系，并区分计划失败、数据写入成功和审计再失败等边界；定向命令 go test ./internal/recovery -run '^TestRestoreReportsPersistenceFailureAfterAuditSucceeds$' -count=1 应稳定证明主错误身份消失，证据基于固定 main SHA，且目标仓库代码、测试和配置零改动。
