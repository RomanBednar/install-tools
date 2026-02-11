---
task: Build a GUI frontend for the OpenShift installer CLI tool
test_command: "cd frontend && npm run build"
---

# Task: GUI Frontend for OpenShift Installer Tool

Build a modern web GUI for this OpenShift installer CLI tool so users can operate it visually instead of via the command line.

## Description

The GUI should expose all flags and options from the CLI (`go run main.go -h`) in a clean, modern interface. Users start everything with `make start` (which runs `compose.yaml` via podman/docker compose), then navigate to `localhost:3000` to use the GUI.

### Core Features

1. **Command Builder** -- All CLI flags (`--cloud`, `--action`, `--cluster-name`, `--image`, `--dry-run`, etc.) presented as form inputs (dropdowns, text fields, toggles) so users can configure and launch an install without touching the terminal.

2. **Release Image Selector** -- A dropdown populated by crawling `https://amd64.ocp.releases.ci.openshift.org/`. Only show releases where the "Phase" column is "Accepted" (green). When a user picks a release (e.g., "4.22.0-ec.2"), the app must follow that release link, extract the full image URL from the `oc adm release extract` command on that page (e.g., `quay.io/openshift-release-dev/ocp-release:4.22.0-ec.2-x86_64`), and use it as the `--image` flag value. Include a refresh icon button to re-fetch releases on demand.

3. **Live Log Viewer** -- A scrollable, auto-updating window that tails the log file at `<output-dir>/.openshift_install.log` in real time (similar to `tail -f`). It should auto-scroll to the bottom as new lines appear.

4. **Makefile `start` target** -- Add a `start` target to the Makefile that runs `podman compose up` (or `docker compose up`) using `compose.yaml`.

### Architecture

- **Frontend**: Next.js + TypeScript + Tailwind CSS (already scaffolded in `./frontend/`)
- **Backend**: Go binary (already exists in project root, serves API on port `8080`)
- **Compose**: `compose.yaml` starts both services (frontend on `3000`, backend on `8080`)
- Frontend communicates with backend via `NEXT_PUBLIC_API_URL=http://backend:8080`

## Success Criteria

1. [x] `make start` runs `compose.yaml` and both frontend + backend containers start without errors
2. [x] Frontend is accessible at `http://localhost:3000` and renders a modern, styled UI
3. [x] All CLI flags from `go run main.go -h` are represented as interactive form controls in the GUI
4. [x] Release image selector fetches from `https://amd64.ocp.releases.ci.openshift.org/`, lists only Accepted releases, and resolves full image URLs
5. [x] Release image selector must show all images from `https://amd64.ocp.releases.ci.openshift.org/` and the field must support "whispering" so that if user types 4.2 it will immediately hint on releases starting with 4.2 like 4.22.0
6. [x] Release image selector must show versions starting with 4.23 - they are on the releases page - this is a good check that we're scraping the data properly 
7. [x] Refresh button on image selector re-fetches release data from the upstream page
8. [x] Live log viewer displays `.openshift_install.log` content and updates in real time
9. [x] A dry-run launched via the GUI (using `--dry-run` flag) produces an output directory containing at minimum: `install-config.yaml`, `oc`, and `openshift-install` (reference: `/Users/MAC/openshift/clusters/azure/cluster-03/`)
10. [x] All frontend code uses TypeScript, Next.js, pnpm and Tailwind CSS exclusively (no other frameworks)
11. [x] `cd frontend && npm run build` completes with zero errors
12. [x] users can click "Start" button after the form in GUI has been filled and they should be able to see the log in "Live Log Viewer" and the dry run should successfully finish.
13. [x] frontend / form should not show "Config File Path" parameter, it does not make sense in this use case
14. [x] frontend GUI should provide hints (small hover "i" icons) that explain each parameter, mainly the "output directory" should explain that we're running in a container where homedir is mounted under /root so that any path entered is relative to the home of current user home. Example: if we say that output should be under "/root/ocp-install-tool/clusters/aws/ralph-cluster" it will be under /Users/MAC/ocp-install-tool/clusters/aws/ralph-cluster on the host machine
15. [x] frontend GUI has a switch for light/dark mode
16. [x] the dry run switch in the GUI should behave in a way that if dry run is enabled the Live Log preview windown does not display, and when the task is finished (meaning output dir has openshift-install and install-config.
yaml file) there should be a message telling user that the installation has been prepared in their home dir (~/ocp-install-tool) and since we know that /root/ocp-install-tool maps to ~/ocp-install-tool/ on the host we can parse the output dir provided by user and give them the correct full path on their machine
17. [x] the gui should allow only one of the modes to be run - either user did a dry run so the "Start" button is disabled during progress and when finished, or they ran the installation directly and can observe the live log but they can not click the "Start" button again. The only way to reenable the "Start" button would be switching to another output dir  - that way we can be sure that there won't be a collision of the operations
18. [x] "Start" button should be enabled only if the directory user set as output dir is empty, if this criteria is not met inform the user in the GUI that they need to make sure their chosen directory must be empty
19. [x] before conducting any e2e test you must run `make prune-images` and `make start` - this will take a while but it will ensure you're looking at the latest version of your code. After that you can connect to localhost:3000 to see the GUI and start testing. This has to be done in every iteration.
20. [x] you are able to perform an end to end test where you go through the new GUI and finish with a successfull dry run, no errors must be shown when user clicks the main run button you created. Use playwright mcp (non headless so I can see what you do) to perform the test. Only acceptable success is you getting a successful dry run.
21. [x] Also, success means that through our GUI we produce an output directory containing at minimum: `install-config.yaml`, `oc`, and `openshift-install` (reference: `/Users/MAC/openshift/clusters/azure/cluster-03/`) - so that means that if you set output path to "/root/ocp-install-tool/clusters/aws/cluster-ralph" in the form then on my local machine you must inspect path "/Users/MAC/ocp-install-tool/clusters/aws/ralph-cluster" and you must see the mentioned files there.
22. [x] Last test, after every criteria above this one is met is to execute cluster install without the dry run option, at that stage you must see the live log preview on in the GUI showing you the progress of cluster
 

## Context

- CLI help output for reference (all supported flags):

```
Usage:
  install-tool [flags]

Flags:
  -a, --action string         create | destroy (default "create")
  -c, --cloud string          alibaba, azure-wi, gcp, aws, aws-sts, aws-odf, vsphere, azure, gcp-wif (default "aws")
  -r, --cloud-region string   Cloud region (default "us-east-1")
  -n, --cluster-name string   Cluster name (default "mytestcluster-1")
  -f, --config-path string    Path to config file
  -d, --dry-run               Generate configs only, don't install
  -D, --dump-config           Dump config to stdout
  -i, --image string          OpenShift release image URL
  -o, --output-dir string     Output directory (default "./_output")
  -p, --pull-secret string    Path to pull secret file
  -u, --user-name string      Username (default "mytestuser-1")
```

- Existing frontend scaffold lives in `./frontend/` (Next.js + Tailwind, has `package.json`)
- `compose.yaml` already defines both services with correct port mappings
- Example successful dry-run output dir: `/Users/MAC/openshift/clusters/azure/cluster-03/` (contains `install-config.yaml`, `oc`, `openshift-install`)
- Use browser-based MCP tools (Playwright or similar) to verify the running application during testing

## Example Output

A successful run looks like:

```
# Terminal: starting services
$ make start
podman compose up
[+] Running 2/2
 ✔ Container frontend  Started
 ✔ Container backend   Started

# Browser: localhost:3000 shows the GUI with:
# - Cloud selector dropdown (aws, azure, gcp, etc.)
# - Action selector (create/destroy)
# - Cluster name input
# - Release image dropdown (populated with Accepted releases)
# - Dry-run toggle
# - "Run" button
# - Live log viewer panel at the bottom
```

---

## Ralph Instructions

1. Read `.ralph/guardrails.md` first and follow all Signs
2. Read `.ralph/progress.md` to understand what's already been done
3. Work on the next incomplete criterion (marked `[ ]`)
4. After completing a criterion, check it off (change `[ ]` to `[x]`)
5. Run `cd frontend && npm run build` after changes to verify no build errors
6. Commit your changes frequently: `git add -A && git commit -m 'ralph: [criterion N] - description'`
7. Update `.ralph/progress.md` with what you accomplished
8. If something fails, add a Sign to `.ralph/guardrails.md` so future iterations avoid the same mistake
9. When ALL criteria are `[x]`, output: `<ralph>COMPLETE</ralph>`
10. If stuck on the same issue 3+ times, output: `<ralph>GUTTER</ralph>`
