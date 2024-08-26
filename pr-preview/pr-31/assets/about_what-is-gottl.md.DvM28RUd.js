import{_ as e,c as t,o as a,a2 as i}from"./chunks/framework.BurqOdzB.js";const b=JSON.parse('{"title":"What is Gottl?","description":"","frontmatter":{},"headers":[],"relativePath":"about/what-is-gottl.md","filePath":"about/what-is-gottl.md"}'),o={name:"about/what-is-gottl.md"},s=i('<h1 id="what-is-gottl" tabindex="-1">What is Gottl? <a class="header-anchor" href="#what-is-gottl" aria-label="Permalink to &quot;What is Gottl?&quot;">​</a></h1><p>This collection serves as a guide to patterns and workflows for building web applications using our preferred tools. It’s designed to help developers create maintainable and scalable applications by providing well-defined workflows that support continuous integration and deployment.</p><h2 id="why-we-built-it" tabindex="-1">Why We Built It <a class="header-anchor" href="#why-we-built-it" aria-label="Permalink to &quot;Why We Built It&quot;">​</a></h2><p>In our previous experience with various ecosystems, we encountered significant complexities that often resulted in codebases filled with hidden abstractions and intricate dependencies. When we transitioned to Go, we were attracted by its simplicity and the absence of opaque, &quot;magic&quot; behaviors. However, as we developed more projects in Go, we observed that our codebases tended to devolve into tangled systems with poorly defined domain boundaries.</p><p>Over time, we found ourselves repeatedly addressing the same architectural challenges—creating new services, managing dependencies, and establishing workflows—often without a consistent approach. This led to redundant effort and inconsistent practices across projects.</p><p>To streamline our development process and establish a foundation for building web applications, we consolidated the patterns and practices that proved effective across various projects. This collection of patterns and workflows is designed to mitigate complexity, enforce clear domain boundaries, and provide a reusable, standardized starting point for Go-based web applications.</p><h2 id="values" tabindex="-1">Values <a class="header-anchor" href="#values" aria-label="Permalink to &quot;Values&quot;">​</a></h2><ul><li><strong>Explicitness</strong>: We value clear, explicit code that is easy to understand and maintain.</li><li><strong>Few tools organized well</strong>: We prefer to use a few well-chosen tools that are organized in a consistent manner.</li><li><strong>Build with Observability in Mind</strong>: We believe that observability is a key aspect of building reliable systems.</li></ul>',8),n=[s];function r(l,d,c,p,h,u){return a(),t("div",null,n)}const f=e(o,[["render",r]]);export{b as __pageData,f as default};