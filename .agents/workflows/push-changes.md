---
description: Automatically create a feature branch (if on main), commit with a generated message, and push to GitHub.
---

1.  **Check Current Branch**:
    ```bash
    git branch --show-current
    ```

2.  **Determine Action**:
    - If the output is `main`:
        - Generate a descriptive feature branch name based on your recent changes (e.g., `feature/add-docker-registry`).
        - Create and switch to the new branch:
          ```bash
          git checkout -b feature/<your-descriptive-name>
          ```
    - If the output is NOT `main`:
        - Proceed on the current branch.

3.  **Stage Changes**:
    - Stage all modified and new files:
      ```bash
      git add .
      ```

4.  **Generate Commit Message**:
    - Create a concise and descriptive commit message that summarizes the work done (e.g., "Add GitHub workflow for Docker publishing and auto-tagging").
    - Commit the changes:
      ```bash
      git commit -m "<your-generated-message>"
      ```

5.  **Push to GitHub**:
    - Push the branch to the remote repository:
      ```bash
      git push origin $(git branch --show-current)
      ```
