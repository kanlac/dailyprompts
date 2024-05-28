# Sidecar

## Q

Sidecar 容器是什么？有什么用？

## A

其定义方式和 initContainer 类似，只不过加上一条 restartPolicy: Always，且经过验证，启动顺序在 initContainer 之后。因此它实际上什么作用，不能解决任何新的问题，只是语义上的、设计模式上的东西。Sidecar 模式帮助我们完成将补充任务从主要业务逻辑中解耦，更好地实现「一个容器干好一件事」的目标。
