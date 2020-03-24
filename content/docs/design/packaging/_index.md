---
title: Packaging
weight: 2
---

# Packaging

The whole idea here is that we want clean package design and project
architecture that is not random. The whole idea around mental models is knowing
where everything is.

If you don't have a Package Oriented Design around your project structure and
your packages, you're really going to have lots of problems.

Packaging is our basic unit of compilation in Go and it's how we structure and
think about applications. The promise is very different from Object Oriented
Design.

I want to try and give you some guidelines and some philosophies so you can,
hopefully, start thinking from an engineering standpoint, how to do this better.
There isn't just one way to do this, I want you to understand, this is the way
that I've developed over the last five years. Every project's different, every
team is different. Dynamics are different. So I don't want you to feel like
you're locked in. I want to make sure is that you've got project structures and
philosophies that are in line, and you've got policies and procedures that are
in line with this.

## Language Mechanics

Here are the language mechanics behind packaging as it relates to Go.
- Packaging directly conflicts with how we've been taught to organize source code.

You've been taught, just create a folder somewhere in your project source tree. Throw some code in there and encapsulate. You know, different parts of your code. And that is because projects like your C++, your Java languages, you're really building a monolithic application as it relates to your source code. In Go, we don't really have a monolithic application. In Go, every folder in your source tree, represents a static library. C allows us to create those A files, those SO files, C# lets us do the DLLs, Java has jar files.

- In other languages, packaging is a feature that you can choose to use or ignore.

Every folder ends up being a static library whether you like it or not. This really puts a constraint on us. So you can't just create a folder randomly and put some stuff in it. You have to think about packaging your component level API design from the very beginning.

- You can think of packaging as applying the idea of microservices on a source tree.

Every folder represents it's own unique API, it's own kind of program boundary. Just like we would have with micro services. So we've got to find ways of decoupling these program boundaries and contracts these program boundaries so they can interact.

- All packages are "first class," and the only hierarchy is what you define in the source tree for your project.

There is no real concept of sub packages, just because a folder sitting inside of another folder doesn't meant that this package is a sub package. So, we want to leverage the hierarchy to give us some indications about relationships, but from the compilers point of view there are no sub packages, there are no relationships, all packages are at the same level.

- There needs to be a way to "open" parts of the package to the outside world.

We talked already, about exporting and unexporting. Exporting being a way of opening up a packages' API to be public.

- Two packages can’t cross-import each other. Imports are a one way street.

It' done to make sure the initialization is consistent. If you have two packages importing each other, which one gets to be initialized first?

It's no more a monolithic application with folders, it is a project with lots of individual static libraries that we call packages that eventually are going to import and be bound together.

## Design Philosophy/Guidelines

I think guidelines of philosophies are very important because they allow you to make sure you're making the right engineering choices, and you're not going off the deep end.

I do believe that the Go community is behind these philosophies. Now, how we implement them, we all have different ways to do it.

- **To be purposeful, packages must provide, not contain.**
  - Packages must be named with the intent to describe what it provides.
  - Packages must not become a dumping ground of disparate concerns.

    Every package must have an API, a clear API on what it provides the user. If you're struggling to name a package, you probably have a code smell there. Now packages like net, fmt, http, os, that's clear. We know what they provide. Packages like util, helper, common, these are packages that contain code, and these are will cause you a lot of problems because they're will be a dumping ground of code, they're will be a single point of dependency, and eventually, you're going to crumble on top. If you have a package today called models, if you have a package today that is just a common set of types, your project has already failed. Types are an artifact to move data across program boundaries. They cannot be an API in and of themselves. Avoid packages like models. Even if you have to duplicate types across multiple packages, that is going to be much better. Remember that concrete data solves the problem. We can leverage interfaces later on to decouple our APIs and to accept everybody's concrete data. You're not just bound to your concrete data through an interface. You can leverage anybody's concrete data. So, I want to really avoid packages like models because now that's a very large thing of containing, and all the sudden, you change one model, one concrete, and you've got cascading changes throughout your code base, very dangerous.

- **To be usable, packages must be designed with the user as their focus.**
  - Packages must be intuitive and simple to use.
  - Packages must respect their impact on resources and performance.
  - Packages must protect the user’s application from cascading changes.
  - Packages must prevent the need for type assertions to the concrete.
  - Packages must reduce, minimize and simplify its code base.
- **To be portable, packages must be designed with reusability in mind.**
  - Packages must aspire for the highest level of portability.
  - Packages must reduce setting policy when it’s reasonable and practical.
  - Packages must not become a single point of dependency.

Portability: the more decoupled, the more reusable a package is, the better it's going to be in your ecosystem. We want to be able to pick up packages and move them. We want to minimize coupling because coupling creates constraints, and policy is a big part of this too. Policy means, is if a package is making decisions about how we log, how we do configuration, how we do things, then only other applications that want to do those things the same can use that package. The more policy a package has, the less reusable it is. We've got to be very clear about when we're setting policy and the purpose of that package will help determine whether we're making the right decisions around policy.

## Package-Oriented Design

I'm gonna share a project structure with you that is my own, that is working, that follows within these design philosophies, and allows us to make some very clear engineering choices across our team, which again, gives us levels of consistency that we need.

### Project Structure

```
Kit                     Application

├── CONTRIBUTORS        ├── cmd/
├── LICENSE             ├── internal/
├── README.md           │   └── platform/
├── cfg/                └── vendor/
├── examples/
├── log/
├── pool/
├── tcp/
├── timezone/
├── udp/
└── web/
```

**Kit**

I really do believe that at least every team, should have a kit project and is bound to a single repo let's say in your GitHub.

A kit project is the set of foundational packages, or APIs, that every application you're building should use.

I would want to see things like log packages, if you're not using the standard library; config packages, if you're not using some third party stuff; web frameworks you're building.

I also want these packages to be as decoupled as possible. In other words, I don't want log importing config or config importing log. These packages shouldn't even be logging at all because if you choose to do some sort of logging, then you're setting policy.

I'm also not a big fan of the logging interface. I mean the standard library doesn't log. If these foundational packages have Goroutines or paths of execution where events are happening that you want to log, then I think you should ask for a handler function and the user can implement whatever they want in the handler function and you just call it during those events. We **don't need interfaces all the time.** Asking for a function, like a handler function, really can help streamline everything and then get you to that logging that you need.

**Application**

Every project that you work on is what I call an application project.

The application project can have multiple binaries in it.

A large team of people can work in an application project if package oriented design is being implemented properly. I've worked on teams of three to five people in this structure and we've been very very successful without stepping on each other.

There are four folders in an application project:
- **vendor/**

    For the purpose of this post, all the source code for 3rd party packages need to be vendored (or copied) into the `vendor/` folder. This includes packages that will be used from the company `Kit` project. Consider packages from the `Kit` project as 3rd party packages.

- **cmd/**

    All the programs this project owns belongs inside the `cmd/` folder. The folders under `cmd/` are always named for each program that will be built. Use the letter `d` at the end of a program folder to denote it as a daemon. Each folder has a matching source code file that contains the `main` package.

- **internal/**

    Packages that need to be imported by multiple programs within the project belong inside the `internal/` folder. One benefit of using the name `internal/` is that the project gets an extra level of protection from the compiler. No package outside of this project can import packages from inside of `internal/`. These packages are therefore internal to this project only.

- **internal/platform/**

    Packages that are foundational but specific to the project belong in the `internal/platform/` folder. These would be packages that provide support for things like databases, authentication or even marshaling.

Example: [Ardan's service project](https://github.com/ardanlabs/service).

Application layer/command layer: packages that are defined inside of let's say "crud" (`cmd/crud/`) are very application specific and they're there to help support start-up, shut-down, maybe some small routing like in a web service. But, not a lot of business logic, more presentational logic, like in a web service. Taking a request, doing a response. You can get away with packages that contain at this level because they have really no level of reusability. They're specific for this app and nothing else. You also could have `cmd/crud/tests/` folders here or packages for integration tests. At the application layer this is where you have a lot of ability to set policy. You could do some containment like we're doing with handlers.

Business service layer logic: `internal/` folder is where we put our business logic. This is going to be business packages, service-level packages, packages that need to be reusable across multiple binaries.

```
├── internal/
│   ├── user/
│   │   └── user.go
│   │   └── user_test.go
│   │   └── models.go
│   ├── mid/
│   │   └── errors.go
│   │   └── logger.go
│   │   └── metrics.go
│   ├── platform/
│   │   ├── db/
│   │   ├── docker/
│   │   ├── flag/
│   │   ├── tests/
│   │   ├── trace/
│   │   ├── web/
```

You could see here that I have a package called "user" that implements all of my user CRUD.

You can also see that I've got a package here called "mid" which is my middleware stuff. That package is providing business layer middleware.

`internal/platform/` is for foundational packages that are not inside of `Kit`. I've got special database support, dockers support, flag support, special testing, tracing, web framework support.

### Validation

**Validate the location of a package.**

- `Kit`
  - Packages that provide foundational support for the different `Application` projects that exist.
  - logging, configuration or web functionality.
- `cmd/`
  - Packages that provide support for a specific program that is being built.
  - startup, shutdown and configuration.
- `internal/`
  - Packages that provide support for the different programs the project owns.
  - CRUD, services or business logic.
- `internal/platform/`
  - Packages that provide internal foundational support for the project.
  - database, authentication or marshaling.

**Validate the dependency choices.**

- `All`
  - Validate the cost/benefit of each dependency.
  - Question imports for the sake of sharing existing types.
  - Question imports to others packages at the same level.
  - If a package wants to import another package at the same level:
    - Question the current design choices of these packages.
    - If reasonable, move the package inside the source tree for the package that wants to import it.
    - Use the source tree to show the dependency relationships.
- `internal/`
  - Packages from these locations CAN’T be imported:
    - `cmd/`
- `internal/platform/`
  - Packages from these locations CAN’T be imported:
    - `cmd/`
    - `internal/`

It's nice to just have one level inside of `platform` and inside of `user`. I think it helps keep things clearer and organized. `platform` is a little harder to try to maintain decoupling at the same level but I work very hard to do this at the `internal` level. I really wanna make sure at the `internal` level, I try to keep packages at the same level from importing each other.

The other thing I want us to do is understand the imports sheet. I have made sure that this projects' structure is in such a way where we should always be able to validate the following.

Here any package that we have inside of command has the ability to import inside of internal and platform. **You can always import down, you can't import up.**
What this means is if I have a package like "user", "user" never ever should be importing packages inside of command.

**Validate the policies being imposed.**

- `Kit`, `internal/platform/`
  - NOT allowed to set policy about any application concerns.
  - NOT allowed to log, but access to trace information must be decoupled.
  - Configuration and runtime changes must be decoupled.
  - Retrieving metric and telemetry values must be decoupled.
- `cmd/`, `internal/`
  - Allowed to set policy about any application concerns.
  - Allowed to log and handle configuration natively.

**Validate how data is accepted/returned.**

- `All`
  - Validate the consistent use of value/pointer semantics for a given type.
  - When using an interface type to accept a value, the focus must be on the behavior that is required and not the value itself.
  - If behavior is not required, use a concrete type.
  - When reasonable, use an existing type before declaring a new one.
  - Question types from dependencies that leak into the exported API.
    - An existing type may no longer be reasonable to use.

**Validate how errors are handled.**

* `All`
    * Handling an error means:
        * The error has been logged.
        * The application is back to 100% integrity.
        * The current error is not reported any longer.
* `Kit`
    * NOT allowed to panic an application.
    * NOT allowed to wrap errors.
    * Return only root cause error values.
* `cmd/`
    * Allowed to panic an application.
    * Wrap errors with context if not being handled.
    * Majority of handling errors happen here.
* `internal/`
    * NOT allowed to panic an application.
    * Wrap errors with context if not being handled.
    * Minority of handling errors happen here.
* `internal/platform/`
    * NOT allowed to panic an application.
    * NOT allowed to wrap errors.
    * Return only root cause error values.

**Validate testing.**

* `cmd/`
    * Allowed to use 3rd party testing packages.
    * Can have a `test` folder for tests.
    * Focus more on integration than unit testing.
* `kit/`, `internal/`, `internal/platform/`
    * Stick to the testing package in go.
    * Test files belong inside the package.
    * Focus more on unit than integration testing.

**Validate recovering panics.**

* `cmd/`
    * Can recover any panic.
    * Only if system can be returned to 100% integrity.
* `kit/`, `internal/`, `internal/platform/`
    * Can not recover from panics unless:
        * Goroutine is owned by the package.
        * Can provide an event to the app about the panic.
