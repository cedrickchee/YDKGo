---
title: Data-oriented design
weight: 6
---

# Data-oriented design

Data-oriented design is an aspect of the language, and one of these things that you're going to have to kind of start switching your brain into away from object-oriented design.

Data-oriented design is about is the understanding that every problem you solve is a data problem. Integrity comes from the data, our performance is going to be coming from the data, everything we do, our mental models, everything's going to be coming from the data. If you don't understand the data you're working with you don't understand the problem you're trying to solve, because all problems are specific and unique to the data that you are working with, and data transformations are at the heart of everything we do, every function, every method, every encapsulation is really around the data transformations that we're working on.

<!--But here's a big part of it, everything we do right now that we've been talking about is in the concrete. Our problems are solved in the concrete data, our manipulations, our memory mutations, everything is in the concrete. But here's the problem. When the concrete data is changing, then your problems are changing, and if your problems are changing, guess what, your algorithms have to change. There's nothing wrong with that except we start to fall into these areas where, once the data's changing and our algorithms change, sometimes that creates changes, cascading changes across an entire code base, and there's a lot of pain.

This is when we start to start focusing on how do I decouple the code from these data changes so these cascading changes are minimized? We will be talking about that in this class when we start getting into methods and interfaces and things like that. But here's the thing, if you're abstracting in a general way, if you're building abstractions on top of abstractions, you're really walking away from all the things we've talked about so far, that the idea of optimizing for correctness and readability. You're walking away from your, really, ability to even understand the mental models of your code because it's too abstract. What we need is this balance of decoupling but thin layers of decoupling to deal with change. Change is still going to be important. But, look, if you're starting to write some code and you're uncertain about the data, this doesn't give you a license to guess. It's a directive to stop, stop what you're doing, sit back and ask yourself, what are the data transformations that are in front of me. You don't have to know about all of them, but you can only code the ones that you're comfortable with, that you feel very confident about in terms of what my input and what my output is.

Look, if you are solving problems you don't have, you're now creating more problems that you do. We are writing code for today, we're going to design an architect for tomorrow. Don't add more code to the code you need today than you do. We've talked about that, that just lends itself to bugs, more lines of code, places for bugs to hide.-->

Data-oriented design is about understanding the data, writing the code that we need, the algorithms that we need and eventually decoupling those algorithms that we have in the concrete to deal with the data changes.

<!--Everything we must do, everything we must do must be focused around minimizing, simplifying, and reducing the amount of code we need to solve every problem.-->

We're going to start looking at those data structures, and start understanding why Go only gives us a arrays, slices, and maps, and I want to show you the mechanical sympathies that these data structures have, and how we're able to not hide the cost of the things that we're doing, and we're able to gain more efficiencies when we're working with the data, and be able to write algorithms that are readable and yet still very, very performant, because, again, on the scale of where do you lose your performance in Go, algorithm efficiency is really at the bottom of that.

<!--It's going to be latency at networking and I/O, latency on allocations and memory, and latencies in your inability to write code that can manage and work with data efficiently.

We're going to start trying to, I'm going to just start trying to break you down and build you back up to get away from these object-oriented design principles and start moving towards these data-oriented design principles.-->
