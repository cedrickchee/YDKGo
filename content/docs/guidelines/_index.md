---
weight: 1
bookFlatSection: true
title: "Design Guidelines"
---

# Design Guidelines

## Prepare Your Mind

- If you look at the Linux operating system as an open-source project, there's about 24 million lines of code in that project today. That is immense; it's massive.
- How many times have you been on a team, and a manager, somebody in management comes and says, "well, you know, "if I throw another developer on this project, "could we get it done faster?" And I'm always like, look, this isn't accounting.
- And there are ways to structure, package projects together where we could put more projects, developers on a team and potentially get more work done.
- A big part of this class is understanding that engineering is about understanding the costs you are taking.
- If you ever met me in a bar and I've been drinking a little whiskey, and you happen to tell me that you're a software developer, well, one of the very first things I will ask you is what is the legacy you're leaving behind?
- And I've got clients today, today that are writing legacy code out of the box, mainly because, one, they don't know how to do it any other way, and two, they create these incredibly short-term milestones for themselves which just really can't be met.
- We've just got to learn how to rethink how we do software development, restructure our software development, and keep this in the front of our minds at all time.
- And I want us to understand that, because all of us at some point will be asked to look at a piece of code that we haven't written.
- But as you continue to write new code, I want you to think about that person that has to maintain that code after you.
- Now if we really are going to focus around not creating legacy code out of the box, for me the number one thing, the number one thing we have to keep in mind are mental models of the code base.
-  A developer, the average developer really can't maintain a mental model of code beyond that 10,000. When we talk about maintaining a mental model, I don't mean memorizing every line of code.
- You understand the code base intimately enough where if there is a problem, you can probably go to the code, identify where to start looking, and start reading that code very, very quickly.
- If the average developer can only handle 10,000 lines of code, then once a project goes beyond 10,000 lines of code, you need a new developer on that project. And when we start talking about code bases with a million lines of code, we're talking about 100 developers on that project if we really want to maintain our mental models and the health of the project.
- You probably don't have even one person that knows the full mental model of that code base.
- Historically, I never allowed developers on my team to use a debugger without asking for permission. This is primarily because when there is a problem in production, all you have is your mental model of the code base and the logs that you're writing.

## Productivity versus Performance

- What I find interesting about our industry, is that for the last 30 years we have kind of been on the same path that Niklaus Wirth has said, back in 1987, he said that "The hope is that the progress in hardware "will cure our software ills."
- Henry, back in 2015, wrote, "The most amazing achievement "of the computer software industry "is its continuing cancellation "of the steady and staggering gains "made by the computer hardware industry." We're still talking about this stuff, 30 years ago, remember I told you before, the hardware is the platform. If performance matters, then we have to be sympathetic with the hardware, mechanical sympathy.
- The hardware folks really haven't changed hardware in the last 15 plus years.
- If we don't understand how Intel 36 core processors works, those processors are actually going to slow us down not speed us up.
- "Go is the best balance I've seen "between power and expressiveness. "You can do almost anything you want to do "by programming fairly straightforwardly, "and have a good mental model "of what's going to happen on the machine, "and you can predict reasonably well, "and how quickly it will run, "and understand what's going on." This is Go in a nutshell.
- When I say mechanics, I'm talking about how things work, when I say semantics, I'm talking about how things behave.

## Correctness versus Performance

- I want to optimize for correctness when we're writing code, I want to use our debuggers and our profiles, not really our debuggers, I apologize, our profiles, our ability to benchmark, our ability to trace, to identify where we may have a performance problem, and then fix the performance issues like that, still trying to maintain levels of correct or readable code.
- The performance of your software will come from four places:
  - latency: on networking and IO, disk IO, those types of things. Those are things you probably want to focus on first.
  - memory allocations: the garbage collector's doing an amazing job, but we don't want to walk away from it.
  - data access: how efficiently you access data.
  - algorithm efficiencies: can we streamline an algorithm to take less steps than more.
- The reality is the hardware will be so fast if we solve the first three problems, that we can have a little algorithm inefficiency, and still get more than we need in terms of performance, and now we're optimizing for correctness, and this is really where we want to go.
- There's a great quote that says that we are one of the few industries where we ask people how to write code before we teach them how to read, and this is really mind boggling when you start to think about it, because everything that we will talk about, about correctness, is going to stem from our ability to read code first.
- I'm trying to bring to you industry leaders who are telling us over and over again correctness over performance when it comes to the code, and to do that, we've gotta be very good at reading code.

## Code Reviews

- People's lives, their happiness, everything is being dependent on the technology that you're building, and we got to take integrity seriously.
- Integrity comes in two places. There's a micro level and a macro level.
- At the micro level, I want you to understand that every line of code you write, every line of code either does one of three things. It allocates memory, reads that memory, and it writes to that memory. What's really interesting is, all our code is doing is reading and writing to memory.
- If you don't understand the data, you do not understand the problem that you're working on, because the problem is the data.
- For every 20 lines of code you write and I write, we've added a bug to the software whether we like it or not.
- Another big part of integrity outside of just writing less code is also dealing with error handling.
- The next step here on our code reviews is readability. Readability means that we're going to structure our software in a way that's comprehensible.
- The average developer on your team should be able to understand the full mental model of the codebase that we're working on, should be able to read every line of code, have a full knowledge of the mental models, and should be able to fix almost any bug, in fact every bug that comes across the codebase.
- What readability also means is not hiding the cost of the code you're writing, not hiding the cost.
- It's very very important that we remember that we are focusing on the real machine, that way we get to be able to write code that doesn't hide the cost, and we give that developer at every moment in time the ability to understand the cost they're taking and the impact it's going to have.
