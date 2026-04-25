---
title: On Software Development
description: My take on the current landscape of software development practices.
created: 2026-04-19
updated: 2026-04-20
---

{{Agentic systems|agentic-systems}} have changed the work of building software. They have changed who can participate, how quickly ideas can become working code, and where the bottlenecks in a team now sit. Professional software teams building products that customers depend on need a deliberate operating model for working with agents in a codebase: one that widens hands-on contribution across product, design, and engineering while keeping clear ownership of the deeper system. When software implementation becomes broadly accessible, teams need clearer definitions of ownership, validation, team structure, and the agent harness they are building around their codebase.

## Agentic Systems Changed the Work

For most of the history of software, engineering held a practical monopoly on changing the product. That was not because every change was strategically deep. It was because even a small change required enough technical fluency to move through the codebase, understand the framework, run the system, and ship safely. If marketing wanted to change the homepage copy, support needed a critical bug fix, or product wanted to try a different user flow through a feature, those requests had to pass through engineering because engineering was the only function that could reliably turn intent into a change.

{{Agentic systems|agentic-systems}} break that bottleneck. They make implementation far more accessible across roles by compressing the distance between an idea and a working change. A product manager can now add analytics instrumentation, a support engineer can trace through a bug and propose a patch, and a marketing team can update presentational UI without waiting on a long queue for engineering time. The change is not just that engineers type faster. The change is that many more people can now participate directly in modifying the software system.

The bottleneck has moved away from manually producing code and toward framing the work, validating the result, and deciding who should own which parts of the system. In other words, the cost of implementation has dropped, but the cost of poor ownership has not. Professional teams now need to rethink their software development lifecycle around that reality.

We are focused here on teams building long-lived software products that customers depend on for business success or meaningful daily use. In that kind of environment, the question is not whether agents are useful. They clearly are. The question is how to adopt them without losing clarity about who understands the system, who is accountable for it, and how it remains maintainable as more people gain the power to change it.

## Software Is Still Information Engineering

Software development is still the work of shaping information systems. We still define the information a system holds (the data), the rules that modify that information, and the ways that information is communicated to users, other systems, and operators. Whether the interface is a web app, an API, a background job, or an event stream, the work is still about managing information correctly through time.

That is why the last 50 years of theory and practice for building information systems still matter. Domain modeling, layered architecture, clear interfaces, stable contracts, thoughtful data modeling, and explicit system boundaries are not artifacts of a pre-agent era. They are responses to the enduring difficulty of building complex systems that remain understandable as they grow. {{Agentic systems|agentic-systems}} may accelerate implementation, but they do not reduce the need for software to be coherent, operable, and easy to change. These three qualities help both humans and agents change the system over time.

This is important because rapid code generation can create the illusion that software development has become primarily about producing working code quickly. It has not. The real work is still deciding how a system should be structured, where responsibilities should live, what abstractions are worth preserving, and how changes in one part of the system should or should not affect another. Code is only the current expression of those decisions.

Professional teams should keep that frame in mind as they adopt agents. If software is still information engineering, then speed alone is not enough. The system still needs to be intelligible. Its behavior still needs to be predictable. Its contracts still need to be clear. And the more people and agents you empower to change it, the more important those underlying structures become.

## Fundamental Limitations of Agents

The next question is whether this new implementation capacity changes who needs to carry judgment and ownership in a software system. I do not think it does. Agents are useful precisely because they can take a goal, work through local context, and produce plausible changes quickly. But that is different from carrying durable responsibility for the system they are changing.

1. Agents are fundamentally oriented around the task in front of them. They can search, edit, run commands, evaluate results, and keep moving toward a local objective. That is valuable. But a complex software system is not just a sequence of local objectives. It is an accumulation of tradeoffs, abstractions, constraints, and organizational decisions that have to remain coherent over time. Agents do not naturally optimize for that larger continuity. They optimize for progress on the current task.

2. Even when an agent performs well within a session, that is not the same thing as owning a system across months or years. Session context is not the same thing as a durable mental model. The long-term understanding of why a system is shaped the way it is, which past decisions still matter, where the risky edges are, and what kinds of changes have historically caused problems still lives primarily in people, teams, documentation, and the structure of the code itself. Agents can help operate inside that context, but they do not naturally become its lasting stewards.

3. This matters even more when a problem calls for a non-obvious solution. The most valuable software decisions usually do not come from producing the most plausible implementation of a prompt. They come from understanding the problem domain deeply enough to notice where the obvious implementation is wrong, incomplete, or strategically unimportant. If development becomes only a matter of prompting from a distance, both engineers and non-engineers can lose contact with the real constraints of the system and the domain it serves.

I do not expect this limitation to disappear as models improve. More capable agents will produce better local work, hold more context, and automate more of the implementation path. But the need for human ownership remains because the core issue is not just raw capability. It is that complex systems still need someone to maintain continuity of judgment across time, across layers, and across changing business needs. Agents can extend that work. They do not eliminate it.

## The Real Risk Is Losing Ownership

Historically, engineering owned implementation because engineering alone had the practical ability to modify the system safely. Even small changes required enough knowledge of the codebase, tooling, deployment model, and failure modes that most other functions had to work through engineering. That constraint was frustrating, but it also created a clear line of ownership. The people changing the system were generally the same people expected to understand it when it failed.

{{Agentic systems|agentic-systems}} change that arrangement. Product, marketing, support, and design can now all contribute more directly to the software itself. In many cases that is a real improvement. It shortens loops, reduces translation loss, and lets the people closest to a problem participate more directly in solving it. But it also creates a new organizational risk: contribution can expand faster than ownership does.

At first glance, it is tempting to conclude that if anyone can change the system, anyone should change the system. This is where teams get into trouble. The burden of understanding incidents, debugging broken flows, protecting data integrity, evolving core abstractions, and keeping the system maintainable does not disappear just because more people can now open pull requests. In most professional environments, that burden still lands on engineering.

If a team adopts broad agentic contribution without a new operating model, engineering can lose practical ownership while still carrying operational accountability. Engineers become responsible for systems they did not shape, patterns they did not approve, and behavior they no longer fully understand. Over time, that is how a team becomes fast at producing changes and slow at restoring trust in the system. The real risk is not that agents write bad code sometimes. The real risk is that teams dissolve ownership boundaries without replacing them with something better.

## Broaden Contribution, Preserve Stewardship

The answer is not to pull non-engineering functions back out of the codebase. The answer is to widen contribution while preserving stewardship. Professional teams should let more people participate directly in improving the product, but they should do so with clear boundaries around which layers are open to broad contribution and which layers require deeper engineering ownership.

In practice, that usually means marketing, product, support, and design can contribute most directly in the client and handler layers of a hypothetical layered architecture. That is where work is closest to presentation, user flows, copy, analytics, and request orchestration. Marketing should not need to wait in an engineering queue to update homepage content, refine a CTA, add a tracking script. Product should be able to add tracking, manage instrumentation, make UI adjustments, and prototype directly in the product instead of only handing engineering a requirements document. Support should be able to reproduce bugs locally, confirm a fix, and open a targeted pull request instead of only filing a ticket and hoping it gets prioritized.

These are meaningful contributions to the information system, but they happen near the edges where the intent is visible and the blast radius is easier to reason about. The deeper layers are different. Service-level contracts, repositories, core business rules, data integrity constraints, and foundational APIs carry more of the system's long-term complexity. Changes there are easier to get subtly wrong and harder to evaluate from the outside. Those layers should remain primarily owned by engineering.

That ownership is not about status. It is about stewardship. Engineering should remain responsible for the deeper contracts that the rest of the system depends on, the structure of the architecture, the reliability of the platform, and the maintainability of the codebase over time. Broader contribution works best when it happens inside a system whose deeper layers are being actively curated. The goal is not to keep people out. The goal is to let more people contribute without dissolving the clarity a professional team needs in order to operate its software well.

## A Better Team Model

The older organizational model for software development was built around handoff. Product wrote a PRD. Design produced mocks and flows. Engineering took those documents, translated them into the system, and carried the implementation to production. That model made sense when the practical ability to change software lived almost entirely inside engineering.

{{Agentic systems|agentic-systems}} make a different model possible. Instead of treating product, design, and engineering as separate phases in a pipeline, teams can work much more collaboratively and much closer to the code itself. Product and design do not have to stop at documents. They can work directly in the user experience, flows, copy, instrumentation, and presentational layers of the product while engineering builds and maintains the deeper building blocks those experiences depend on.

In practice, that means one or two engineers can own a feature end-to-end in close partnership with product and design. Those engineers are responsible for the technical shape of the feature: the underlying data models, core services, APIs, system boundaries, validation, and operational readiness. Product and design work directly with them, often hands-on in the codebase, consuming those building blocks in the interface and iterating on the user journey as the feature takes shape.

This is not a model where everyone owns everything. It is a model where collaboration becomes tighter while technical responsibility stays clear. Engineers own the durable building blocks and deeper contracts that the rest of the product depends on. Product and design work more directly in the layers that express user intent and consume those building blocks. A senior, staff, or principal engineer helps steward the architecture of the whole system across features so that local speed does not erode global coherence.

That kind of team can move much faster with agents now and in the long term, but only if the work is carried through to completion. A feature is not done when the UI appears to work. It is done when the underlying contracts are sound, the behavior is validated, the operational edges are understood, and the agent harness has been updated so the next contributor can extend the work safely.

## Build the Harness

If professional teams want that kind of collaborative model to work, they have to invest directly in the agent harness around the code. An agent harness is not the model itself, and it is not just an engineer or product manager using an agent in a repository from time to time. It is the team-specific system built around agents, local environments, repository guidance, validation, automation, and review. It is what turns raw model capability into a repeatable way of building a particular software product.

That distinction matters because agents on their own are generic. They can generate, edit, search, and iterate, but they do not automatically know how your product works, which architectural boundaries matter, what your failure modes look like, which patterns are acceptable, or how your team wants changes to be validated before they enter the system. If more people are going to contribute directly to a software system, they need more than access to a model. They need an environment that helps them work against reality, guidance that helps them respect the system's shape, and validation that helps them know when a change is good enough to hand off, review, or merge.

Every serious team will end up building a different harness because every serious team has a different system. The right harness for a small product with a simple operational profile is not the same as the right harness for a multi-tenant SaaS platform with complex workflows and strict uptime expectations. This is not optional process overhead. It is part of the product development system now. The easier it becomes to produce code, the more important it becomes to shape the environment that code is produced inside.

### 1. Production-like local environments

Contributors and agents need a local development environment that closely mirrors production behavior. If they cannot run the real application, exercise meaningful workflows, inspect system behavior, and verify outcomes locally, they cannot truly finish tasks. They can only produce plausible-looking patches. In practice, that means local infrastructure, realistic data flows, and development setups that expose the same important boundaries the production system does. When production must diverge from local functionality, those differences should sit behind thin interfaces with local-only implementations so the rest of the system still behaves consistently.

### 2. Progressive guidance inside the work

A good harness meets contributors inside the work itself. It progressively exposes architecture, patterns, anti-patterns, documentation, and {{cross-cutting concerns|cross-cutting-concerns}} while changes are being made. It does not assume that every contributor, or every agent, already understands the system deeply. Instead, it makes important constraints visible at the moment they matter. The better the harness teaches the shape of the system during implementation, the less review energy gets wasted correcting avoidable mistakes later.

### 3. Enforcement through automation

Formatting, linting, builds, type checks, structural checks, and review workflows should all make the right path the easiest path. If a team has preferred module structures, dependency directions, testing seams, or interface boundaries, the harness should encode those expectations. A professional software team should not rely on every contributor to memorize the system perfectly. It should put as much of that knowledge as possible into the development environment itself.

## A Validation System, Not Just More Tests

Fast-moving software needs validation at the right levels of abstraction, not just more tests for their own sake. The point of a validation system is to create durable confidence as implementation changes rapidly underneath it. That matters even more in a collaborative model where product, design, and engineering can all help shape a feature while engineers remain accountable for the deeper system. If the checks are too close to the current shape of the code, they become expensive to maintain and too easy to bypass mentally. If they are too shallow, they do not tell you whether the system actually works.

### 1. Prefer handler-level tests

My default preference is to validate behavior at the handler level. Send a real HTTP request, provide a realistic payload, assert the response, and verify the important side effects in the system. That might mean confirming a database update, checking that an email was dispatched, or asserting that a domain event was emitted. These tests sit at a high enough level that the implementation underneath can change without forcing the test to change with it. That makes them a strong fit for a codebase where both humans and agents are going to be reshaping internals frequently.

### 2. Use as few mocks as possible

I want as few mocks as possible in that flow. In particular, I do not want to mock the database unless there is a very specific reason to do so. The database is too central to most backend systems to treat as an optional detail during validation. External APIs are different. Those should usually be hidden behind clean interfaces and mocked at that seam, ideally with type-safe dependency boundaries. That gives the system a realistic core while still letting tests stay deterministic around third-party integrations.

### 3. Add focused service-level tests

Handler-level tests are not enough on their own. Some important service boundaries need focused tests, especially where failure modes are difficult to drive from a higher-level request or where error handling carries real business significance. Those tests should exist to cover meaningful gaps, not to duplicate what is already proven elsewhere. If a service has important retry behavior, transaction semantics, or subtle failure handling, that is a good place for a more direct test.

### 4. Be selective on the client

I am much more selective about client-side testing. Client code is event-driven, easy to reshuffle, and often expensive to test in a way that stays useful over time. In many cases I would rather regenerate, simplify, or rewrite complex client code than accumulate a large fragile test suite around it. That does not mean never test the client. It means the bar should be higher. Validation effort should go first to the places where system behavior, persistence, and contracts actually create the most risk.

## Code Review and Refactoring Matter More, Not Less

In an agentic workflow, senior engineers spend less time manually typing every change and more time governing what enters the system. That does not make review less important. It makes review more important, because the volume of code can increase much faster than the team's shared understanding of that code. Review becomes one of the main places where a team protects the integrity of the system and keeps the boundary between core building blocks and feature-level consumption clear.

That means code review should change in emphasis. The goal is not to bikeshed every naming choice or push every submission toward some idealized notion of handcrafted elegance. Most of the time, mediocre code is fine if it is safe, legible, and easy to improve later. The review bar should focus on harmful patterns, unclear ownership, broken boundaries, missing validation, operational risk, and anything likely to make the system harder for the next engineer, product partner, designer, or agent to understand.

The deeper goal is not just to reject bad changes. It is to improve the agent harness so the same class of bad change becomes less likely to happen again. If a review repeatedly catches missing instrumentation, broken layering, weak tests, or inconsistent module structure, the team should encode that lesson into guidance, tooling, or automation. In a healthy system, review is one feedback loop in a larger process of making the harness better.

That connects directly to refactoring. {{Agentic systems|agentic-systems}} can increase code volume much faster than they increase coherence. Left alone, they will slowly rot a codebase by adding duplication, fuzzy boundaries, and layers of accidental complexity. Professional teams should treat continuous refactoring as core work, not leftover work. Some portion of the engineering organization should be dedicated to cleanup, simplification, and structural maintenance so the product's shared building blocks stay legible and both humans and agents can continue to understand and modify the system effectively.

## Contain Vibe Coding

I think it is useful to separate {{vibe coding|vibe-coding}} from the rest of this discussion. Vibe coding is a real mode of working, and it can be productive. You are not trying to understand or shape the code deeply. You are iterating on the experience the code produces. If the output works, feels right, and is cheap to discard, that is often enough.

That mode has a place. It can be excellent for prototypes, internal tools, throwaway experiments, and isolated leaf nodes where the blast radius is small and the long-term maintenance burden is low. In those cases, speed of exploration matters more than code quality, and the generated artifact may only need to survive long enough to prove a point.

But professional software systems cannot be built entirely that way. Once software carries uptime expectations, data integrity requirements, operational obligations, and long-term product ownership, the code itself matters again. Someone has to read it, reason about it, change it, and trust it under pressure. That is where vibe coding stops being a complete development model.

So I do not think professional teams should reject vibe coding. I think they should contain it. Use it where the system can tolerate disposability and weak ownership. Be much more disciplined where the software carries lasting responsibility. In practice, that usually means vibe coding belongs in exploratory flows, prototypes, and leaf features at the edges of the system, not in the core building blocks that other contributors depend on.

## Adopt the Change Professionally

Professional teams should adopt {{agentic systems|agentic-systems}}. The change is real, the leverage is real, and the opportunity to widen participation in software development is real. Teams that ignore that will give up both speed and organizational learning.

But adoption is not just a matter of giving everyone access to an agent and hoping for the best. The real work is replacing a handoff-driven model with a collaborative one that still preserves clear technical ownership. That means product, design, and engineering working closer together in the code, engineers owning the core building blocks and deeper contracts, and senior engineers maintaining architectural coherence across the system.

That model only works when teams invest in the structures around it. They need stronger validation, better local environments, more intentional review, continuous refactoring, and a deliberate agent harness around the codebase. They also need to treat a feature as finished only when the contracts are sound, the behavior is validated, and the harness has enough guidance for the next human or agent to extend the work safely.

If agentic systems changed the work, then professional teams should change the way they work in response. Let small teams own features end-to-end. Let product and design work hands-on in the layers closest to user intent. Let engineers maintain the building blocks the product depends on. And build the harness that makes that model durable.
