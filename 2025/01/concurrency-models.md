# 并发模型（架构）

## 线程驱动与事件驱动及各自优势

线程驱动：编程简单直接，不需要处理异步和回调。例：Tomcat，Apache web server，每个请求由一个线程处理；GMP 模型

事件驱动：通常使用单线程事件循环或少量线程处理事件队列（event queue），减少了创建线程的开销，也减少了竞态问题，适合高并发。例：Nginx。

## 多进程并发 v.s. 多线程并发

### 多进程模型

子进程可以并行运行在多个处理器上，也可以被同一个处理器通过上下文切换进行 CPU 分时共享。创建进程和上下文切换的开销较大。

在涉及共享资源时，可能会有竞态问题。

### 多线程模型

通过 clone 系统调用创建线程。线程与进程的区别在于同一个进程内的线程**共享**着进程分配的资源，线程不被分配资源，只是操作系统调度执行任务的抽象的最小单元。

因为所有线程共享同一个进程的内存空间，所以**更容易**出现竞态问题。

### 实例

- Apache web server 是采用多进程与多线程混用的模型，在每个用户请求接入的时候都会通过 fork 系统调用创建一个子进程，如此同时支持多个用户，然后每个进程会派生线程，每个线程处理一条连接。
- PostgreSQL 使用多进程架构，每个连接一个进程
