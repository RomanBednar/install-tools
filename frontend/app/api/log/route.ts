import { NextRequest, NextResponse } from "next/server";

export const dynamic = "force-dynamic";

const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export async function GET(request: NextRequest) {
  try {
    const path = request.nextUrl.searchParams.get("path");
    const url = path
      ? `${BACKEND_URL}/log?path=${encodeURIComponent(path)}`
      : `${BACKEND_URL}/log`;

    const response = await fetch(url, {
      cache: "no-store",
    });

    if (!response.ok) {
      // Return empty log if not available yet
      if (response.status === 500) {
        return NextResponse.json({ log: "" });
      }
      return NextResponse.json(
        { error: `Failed to fetch log: ${response.status}` },
        { status: response.status }
      );
    }

    const logContent = await response.text();
    return NextResponse.json({ log: logContent });
  } catch (error) {
    // Return empty log on error (log file may not exist yet)
    return NextResponse.json({ log: "" });
  }
}
