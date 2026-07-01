# web-apps

Tiny tenant application used to test PR preview environments.

The important contract is the image tag:

```text
ghcr.io/notniknot/web-apps:preview-pr-<number>
```

The GitHub Actions workflow publishes this tag for pull requests. Preview
automation then deploys that mutable tag and lets Argo CD Image Updater pin the
current digest back into `notniknot/web-apps-config`.

