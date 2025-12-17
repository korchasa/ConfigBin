---
description:
globs:
alwaysApply: true
---
## DOCUMENTATION STRUCTURE AND RULES (directory `documents`)

### Hierarchy and purpose

- Software Requirements Specification (SRS): Is the primary source of truth for the project. Answers the questions: what are we doing and why.
- Software Design Specification (SDS): Is a source of project implementation details. Depends on Software Requirements Specification (SRS). Answers the question: how we do it.
- File Structure Map: A map of the project's file structure and its purpose.
- Whiteboard: A temporary notes file for in-progress notes.

### Documentation Rules

- Application MUST STRICTLY COMPLY with the SRS and SDS.
- When adding a new requirement or updating existing ones: update SRS -> update SDS -> implementation.
- Implemented requirements and acceptance criteria should be marked with [x] before the requirement and criterion title. Not implemented ones should be marked with [ ] or omitted.

### Software Requirements Specification (SRS) Format (file @documents/requirements.md)

```markdown
# Software Requirements Specification (SRS)

## 1. Introduction
- **Document purpose:**
- **Scope:**
- **Audience:**
- **Definitions and abbreviations:**

## 2. General description
- **System context:** (diagram or environment description)
- **Assumptions and constraints:**
- **Assumptions:**

## 3. Functional requirements
### 3.1 Requirement FR-1
- **Description:**
- **Use case scenario:**
- **Acceptance criteria:**

### 3.2 Requirement FR-2
...

## 4. Non-functional requirements
- **Performance:**
- **Reliability:**
- **Security:**
- **Scalability:**
- **Availability/UX:**

## 5. Interfaces
- **APIs and integrations:**
- **Protocols and data formats:**
- **UI/UX constraints:**

## 6. Acceptance criteria
- The system is considered accepted if the following are met: ...
```

### Software Design Specification (SDS) Format (file @documents/design.md)

```markdown
# Software Design Specification (SDS)

## 1. Introduction
- **Document purpose:**
- **Relation to SRS:** (links to requirements)

## 2. System Architecture
- **Overview diagram:** (C4/UML/block diagram)
- **Main subsystems and their roles:**

## 3. Components
### 3.1 Component A
- **Purpose:**
- **Interfaces:** (API, input/output)
- **Dependencies:**

### 3.2 Component B
...

## 4. Data and Storage
- **Entities and attributes:**
- **ER diagram:**
- **Migration policies:**

## 5. Algorithms and Logic
- **Key algorithms:** (pseudocode or diagram)
- **Business rules:**

## 6. Non-functional Aspects
- **Scalability:**
- **Fault tolerance:**
- **Security:**
- **Monitoring and logging:**

## 7. Constraints and Trade-offs
- What has been simplified
- What has been deferred to future versions

## 8. Future Extensions
- Ideas and opportunities for the roadmap
```

### File Structure Map Format (file @documents/file_structure.md)

Parts of the file:
- tree-view tree(with file purposes and relationships) for:
  + root directory
  + sources
  + tests
- file organization patterns
- English language only

### Whiteboard Format (file @documents/whiteboard.md)

- Temporary notes
- Ongoing plans and progress marks.
- The only file for in-progress notes.
- Must be cleaned up after new session starts.



