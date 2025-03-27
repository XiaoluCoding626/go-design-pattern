# go-design-pattern

使用 Go 语言完整实现 58 种设计模式，包含 GoF 23 种经典模式与现代设计模式的详细示例与最佳实践。

## 项目简介

本项目旨在通过 Go 语言实现各种设计模式，帮助开发者深入理解和应用设计模式。每种模式都包含详细的说明、UML类图、代码实现、使用示例和最佳实践，便于学习和参考。

## 设计模式分类

### 创建型模式 (Creational Patterns)

这些模式与对象的创建机制有关，主要解决对象创建时的灵活性和复用性问题。

- [x] [简单工厂模式（Simple Factory）](./creational/simple_factory/docs/README.md)
- [x] [工厂方法模式 (Factory Method)](./creational/factory_method/docs/README.md)
- [x] [抽象工厂模式 (Abstract Factory)](./creational/abstract_factory/docs/README.md)
- [x] [建造者模式 (Builder)](./creational/builder/docs/README.md)
- [x] [原型模式 (Prototype)](./creational/prototype/docs/README.md)
- [x] [单例模式 (Singleton)](./creational/singleton/docs/README.md)
- [x] [New 模式 (New)](./creational/new/docs/README.md)
- [x] [函数选项模式 (Functional Options)](./creational/functional_options/docs/README.md)
- [x] [对象池模式 (Object Pool)](./creational/object_pool/docs/README.md)

### 行为型模式 (Behavioral Patterns)

这些模式关注对象之间的责任分配和算法封装，以及对象间的通信方式。

- [x] [策略模式 (Strategy)](./behavioral/strategy/docs/README.md)
- [x] [命令模式 (Command)](./behavioral/command/docs/README.md)
- [x] [观察者模式 (Observer)](./behavioral/observer/docs/README.md)
- [x] [访问者模式 (Visitor)](./behavioral/visitor/docs/README.md)
- [x] [迭代器模式 (Iterator)](./behavioral/iterator/docs/README.md)
- [x] [模板方法模式 (Template Method)](./behavioral/template_method/docs/README.md)
- [x] [状态模式 (State)](./behavioral/state/docs/README.md)
- [x] [备忘录模式 (Memento)](./behavioral/memento/docs/README.md)
- [x] [中介者模式 (Mediator)](./behavioral/mediator/docs/README.md)
- [x] [解释器模式 (Interpreter)](./behavioral/interpreter/docs/README.md)
- [x] [责任链模式 (Chain of Responsibility)](./behavioral/chain_of_responsibility/docs/README.md)
- [x] [注册表模式（Registry）](./behavioral/registry/docs/README.md)
- [x] [上下文模式（Context）](./behavioral/context/docs/README.md)

### 结构型模式 (Structural Patterns)

这些模式关注类和对象的组合，形成更大的结构，同时保持结构的灵活和高效。

- [x] [适配器模式 (Adapter)](./structural/adapter/docs/README.md)
- [x] [桥接模式 (Bridge)](./structural/bridge/docs/README.md)
- [x] [组合模式 (Composite)](./structural/composite/docs/README.md)
- [x] [装饰器模式 (Decorator)](./structural/decorator/docs/README.md)
- [x] [外观模式 (Facade)](./structural/facade/docs/README.md)
- [x] [享元模式 (Flyweight)](./structural/flyweight/docs/README.md)
- [x] [代理模式 (Proxy)](./structural/proxy/docs/README.md)

### 同步模式（Synchronization Patterns）

这些模式用于处理并发编程中的同步问题。

- [x] [互斥锁模式 (Lock Mutex)](./synchronization/lock_mutex/docs/README.md)
- [x] [读写锁模式 (Read-Write Lock)](./synchronization/read_write_lock/docs/README.md)
- [x] [信号量模式 (Semaphore)](./synchronization/semaphore/docs/README.md)
- [x] [条件变量模式 (Condition Variable)](./synchronization/condition_variable/docs/README.md)
- [x] [监视器模式 (Monitor)](./synchronization/monitor/docs/README.md)

### 并发模式 (Concurrency Patterns)

这些模式用于处理并发编程中的各种问题。

- [ ] 生产者-消费者模式 (Producer-Consumer)
- [ ] 读写锁模式 (Read-Write Lock)
- [ ] 线程池模式 (Thread Pool)

### 架构模式 (Architectural Patterns)

这些模式用于系统架构的设计。

- [ ] MVC 模式 (Model-View-Controller)
- [ ] MVVM 模式 (Model-View-ViewModel)
- [ ] 微服务架构 (Microservices)

## 项目结构
