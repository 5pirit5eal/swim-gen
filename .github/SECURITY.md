# Security Policy

## Reporting a Vulnerability

Report suspected security vulnerabilities privately to `info.swim.gen@gmail.com`.
Please do not open a public GitHub issue for an unremediated vulnerability.

Include the affected URL or component, a concise description of the issue, steps
to reproduce it, the potential impact, and any evidence that can help reproduce
the behavior. Remove personal data, credentials, access tokens, and other
sensitive information from the report whenever possible.

Swim Gen will acknowledge reports when practical and will coordinate any public
disclosure with the reporter. This project does not currently offer a bug bounty
or guaranteed response time.

## Scope

The primary scope is the production application at `https://swim-gen.com` and
the source code in this repository. Testing must use accounts and data that you
own or are explicitly authorized to use.

Supabase, Google Cloud, and other third-party services have their own security
reporting processes. Report vulnerabilities in those services to their
respective providers.

## Good-Faith Testing

Good-faith testing is welcome when it avoids privacy violations, data access
outside the tester's own accounts, denial-of-service activity, destructive
changes, social engineering, phishing, or attacks against third parties.

Stop testing and report the issue privately if testing could affect another
person's data or service availability.
