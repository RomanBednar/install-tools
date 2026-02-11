import { NextResponse } from "next/server";

export const dynamic = "force-dynamic";

interface ReleaseTag {
  name: string;
  phase: string;
  pullSpec: string;
  downloadURL?: string;
}

interface StreamResponse {
  name: string;
  tags: ReleaseTag[];
}

interface Release {
  name: string;
  phase: string;
  imageUrl: string;
  stream: string;
}

// All release streams shown on https://amd64.ocp.releases.ci.openshift.org/
const RELEASE_STREAMS = [
  "4-stable",
  "4-dev-preview",
  "5-stable",
  "5.0.0-0.nightly",
  "4.23.0-0.nightly",
  "4.23.0-0.ci",
  "4.22.0-0.nightly",
  "4.22.0-0.ci",
  "4.21.0-0.nightly",
  "4.21.0-0.ci",
  "4.20.0-0.nightly",
  "4.20.0-0.ci",
  "4.19.0-0.nightly",
  "4.19.0-0.ci",
  "4.18.0-0.nightly",
  "4.18.0-0.ci",
  "4.17.0-0.nightly",
  "4.17.0-0.ci",
  "4.16.0-0.nightly",
  "4.16.0-0.ci",
];

const BASE_URL = "https://amd64.ocp.releases.ci.openshift.org";

async function fetchStream(stream: string): Promise<Release[]> {
  try {
    const resp = await fetch(`${BASE_URL}/api/v1/releasestream/${stream}/tags`, {
      headers: { Accept: "application/json" },
      next: { revalidate: 0 },
    });

    if (!resp.ok) return [];

    const data: StreamResponse = await resp.json();
    if (!data.tags) return [];

    return data.tags
      .filter((tag) => tag.phase === "Accepted")
      .map((tag) => ({
        name: tag.name,
        phase: tag.phase,
        imageUrl: tag.pullSpec,
        stream: stream,
      }));
  } catch {
    return [];
  }
}

export async function GET() {
  try {
    // Fetch all streams in parallel
    const results = await Promise.all(RELEASE_STREAMS.map(fetchStream));
    const allReleases = results.flat();

    // Deduplicate by name (same version might appear in multiple streams)
    const seen = new Set<string>();
    const releases: Release[] = [];
    for (const release of allReleases) {
      if (!seen.has(release.name)) {
        seen.add(release.name);
        releases.push(release);
      }
    }

    // Sort by version descending
    releases.sort((a, b) => compareVersions(b.name, a.name));

    return NextResponse.json({ releases });
  } catch (error) {
    console.error("Error fetching releases:", error);
    return NextResponse.json(
      { error: `Failed to fetch releases: ${error}` },
      { status: 500 }
    );
  }
}

function compareVersions(a: string, b: string): number {
  const parseVersion = (v: string) => {
    const match = v.match(/^(\d+)\.(\d+)\.(\d+)/);
    if (!match) return [0, 0, 0];
    return [parseInt(match[1]), parseInt(match[2]), parseInt(match[3])];
  };

  const va = parseVersion(a);
  const vb = parseVersion(b);

  for (let i = 0; i < 3; i++) {
    if (va[i] !== vb[i]) return va[i] - vb[i];
  }

  return a.localeCompare(b);
}
