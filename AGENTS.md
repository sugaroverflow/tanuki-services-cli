# Agent Instructions

## GitHub Actions to GitLab CI/CD migrations

When a user asks to migrate, translate, convert, review, or compare GitHub Actions workflows with GitLab CI/CD, load and follow the `github-actions-to-gitlab-ci` skill in full.

Before giving the final response, complete the full skill deliverables. Treat this as a required delivery gate, not optional guidance.

Final response checklist:

- [ ] Review the project context and source GitHub Actions workflows.
- [ ] Create, update, or show the translated `.gitlab-ci.yml`.
- [ ] Validate the translated GitLab CI/CD YAML with `glab ci lint`, the skill validator, GitLab CI Lint, or local YAML validation. If validation cannot be run, state that plainly.
- [ ] Include the validation result in the response.
- [ ] Include a migration report table with: source item, category, GitLab equivalent, and notes.
- [ ] Include manual setup steps.
- [ ] Include a suggested GitLab improvements summary, even if no major improvements apply.
- [ ] Include honest caveats and items that need human review.

Do not end the migration task after only creating or validating `.gitlab-ci.yml`. The final response must include the full handoff package above.
