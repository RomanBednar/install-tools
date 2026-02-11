import { NextRequest, NextResponse } from "next/server";

export const dynamic = "force-dynamic";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const {
      username,
      sshPublicKeyFile,
      pullSecretFile,
      outputDir,
      clusterName,
      image,
      cloudRegion,
      cloud,
      dryRun,
      action,
    } = body;

    // Step 1: Save config to backend
    const saveResponse = await fetch(`${BACKEND_URL}/save`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        username: username || "",
        sshPublicKeyFile: sshPublicKeyFile || "",
        pullSecretFile: pullSecretFile || "",
        outputDir: outputDir || "",
        clusterName: clusterName || "",
        image: image || "",
        cloudRegion: cloudRegion || "",
        cloud: cloud || "",
        dryRun: dryRun ? "true" : "false",
        action: action || "create",
      }),
    });

    if (!saveResponse.ok) {
      const errorText = await saveResponse.text();
      return NextResponse.json(
        { error: `Failed to save config: ${errorText}` },
        { status: saveResponse.status }
      );
    }

    // Step 2: Trigger the action
    const actionResponse = await fetch(`${BACKEND_URL}/action`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ action: action || "create" }),
    });

    if (!actionResponse.ok) {
      const errorText = await actionResponse.text();
      return NextResponse.json(
        { error: `Failed to start action: ${errorText}` },
        { status: actionResponse.status }
      );
    }

    const result = await actionResponse.json();
    return NextResponse.json(result);
  } catch (error) {
    console.error("Error in run endpoint:", error);
    return NextResponse.json(
      { error: `Internal error: ${error}` },
      { status: 500 }
    );
  }
}
