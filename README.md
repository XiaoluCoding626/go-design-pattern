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
- [ ] 函数选项模式 (Functional Options)
- [ ] 对象池模式 (Object Pool)

### 结构型模式 (Structural Patterns)

这些模式关注类和对象的组合，形成更大的结构，同时保持结构的灵活和高效。

- [x] [适配器模式 (Adapter)](./structural/adapter/docs/README.md)
- [ ] 桥接模式 (Bridge)
- [ ] 组合模式 (Composite)
- [ ] 装饰器模式 (Decorator)
- [ ] 外观模式 (Facade)
- [ ] 享元模式 (Flyweight)
- [ ] 代理模式 (Proxy)

### 行为型模式 (Behavioral Patterns)

这些模式关注对象之间的责任分配和算法封装，以及对象间的通信方式。

- [x] [策略模式 (Strategy)](./behavioral/strategy/README.md)
- [ ] 命令模式 (Command)
- [ ] 观察者模式 (Observer)
- [ ] 访问者模式 (Visitor)
- [ ] 迭代器模式 (Iterator)
- [ ] 模板方法模式 (Template Method)
- [ ] 状态模式 (State)
- [ ] 备忘录模式 (Memento)
- [ ] 中介者模式 (Mediator)
- [ ] 解释器模式 (Interpreter)
- [ ] 责任链模式 (Chain of Responsibility)

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
