---
---

# What is Gottl?

This collection serves as a  guide to patterns and workflows for building web applications using our preferred tools. It’s designed to help developers create maintainable and scalable applications by providing well-defined workflows that support continuous integration and deployment.

## Why We Built It

In our previous experience with various ecosystems, we encountered significant complexities that often resulted in codebases filled with hidden abstractions and intricate dependencies. When we transitioned to Go, we were attracted by its simplicity and the absence of opaque, "magic" behaviors. However, as we developed more projects in Go, we observed that our codebases tended to devolve into tangled systems with poorly defined domain boundaries.

Over time, we found ourselves repeatedly addressing the same architectural challenges—creating new services, managing dependencies, and establishing workflows—often without a consistent approach. This led to redundant effort and inconsistent practices across projects.

To streamline our development process and establish a foundation for building web applications, we consolidated the patterns and practices that proved effective across various projects. This collection of patterns and workflows is designed to mitigate complexity, enforce clear domain boundaries, and provide a reusable, standardized starting point for Go-based web applications.

## Values

- **Explicitness**: We value clear, explicit code that is easy to understand and maintain.
- **Few tools organized well**: We prefer to use a few well-chosen tools that are organized in a consistent manner.
- **Build with Observability in Mind**: We believe that observability is a key aspect of building reliable systems.