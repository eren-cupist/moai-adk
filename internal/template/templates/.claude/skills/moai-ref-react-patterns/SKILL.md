---
name: moai-ref-react-patterns
description: >
  Review rubric for React and Next.js components: where state lives, data fetching,
  client/server boundaries, component structure, rendering cost and accessibility. Used
  when implementing UI and by the /moai review UX perspective. Not for framework version
  facts or design direction (moai-domain-frontend), or backend work.

when_to_use: >
  Load when writing or reviewing React components, hooks or client state.

user-invocable: false
metadata:
  version: "2.0.0"
  category: "domain"
  status: "active"
  updated: "2026-03-30"
  tags: "react, nextjs, components, state, accessibility, review, reference"

# MoAI Extension: Progressive Disclosure
progressive_disclosure:
  enabled: true
  level1_tokens: 100
  level2_tokens: 1000
---

# React Component Rubric

Follow the project's existing patterns first: its folder layout, its state and query libraries, its component kit. Report each item below that a change violates as a finding with severity and confidence.

## Where state lives

| State | Home |
|-------|------|
| Restorable from the URL (filters, tabs, pagination) | search params |
| Server data | server components, or the project's query library on the client |
| Form input and validation | the project's form library, with the same schema validated on the server |
| Shared client-only UI state | the project's store, or context for rarely changing values |
| Local to one component | `useState` or `useReducer` |

Findings: server data copied into a client store or `useState`; state derivable from props or other state stored separately and synced with an effect; URL-shaped state kept only in memory so reload or share loses it.

## Data and boundaries

- No `useEffect` plus `fetch` for data a server component or the query library could load; no request waterfalls where requests could run in parallel.
- `'use client'` sits on the smallest subtree that needs interactivity; server-only modules (database clients, secrets) are never imported into client code.
- Every async view has loading, empty and error states, and an error boundary (`error.tsx` or equivalent) covers the route.
- Mutations show pending state, handle failure visibly, and refresh or revalidate the data they change.

## Structure

- Props are typed; no `any` in props or state.
- A component that mixes data loading, business rules and markup, or has grown past what fits on a screen, is split along those lines.
- Reused logic lives in a hook; reused markup in a component from the shared kit.
- Effects have correct dependencies and clean up subscriptions and timers; effects are not used to respond to events that a handler could handle.

## Rendering cost

- Lists render with stable keys (not array indexes when items reorder); long lists are virtualized.
- Large client dependencies are loaded lazily; images use the framework's image component with sizes set.
- Memoization is added where a profile shows a cost, not by default (and not at all when the React Compiler handles it).

## Accessibility

- Interactive elements are real `button`, `a` and form elements, reachable and operable by keyboard, with a visible focus style.
- Every control has an accessible name; images have meaningful `alt` (empty for decorative ones).
- Dialogs and menus trap and restore focus and close on Escape.
- Text and essential UI meet WCAG AA contrast; information is not conveyed by color alone; motion respects `prefers-reduced-motion`.
